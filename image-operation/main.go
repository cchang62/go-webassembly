package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"math"
	"runtime"
	"syscall/js"
	"time"

	"github.com/disintegration/imaging"
)

func Resize(this js.Value, args []js.Value) any {
	// retrieve args
	w, h := args[0].Int(), args[1].Int()
	imgBytesLen := args[2].Get("length").Int()
	imgBytes := make([]byte, imgBytesLen)
	js.CopyBytesToGo(imgBytes, args[2])

	// process
	src, _, _ := image.Decode(bytes.NewReader(imgBytes))
	newImg := imaging.Resize(src, w, h, imaging.Lanczos)

	// output
	buffer := new(bytes.Buffer)
	imaging.Encode(buffer, newImg, imaging.JPEG)
	switch args[3].String() {
	case "string":
		return base64.StdEncoding.EncodeToString(buffer.Bytes())
	default:
		jsArray := js.Global().Get("Uint8Array").New(args[2].Get("length").Int())
		js.CopyBytesToJS(jsArray, buffer.Bytes())
		return jsArray
	}
}

func FlipImageH(this js.Value, args []js.Value) any {
	// retrieve args
	imgBytesLen := args[0].Get("length").Int()
	imgBytes := make([]byte, imgBytesLen)
	js.CopyBytesToGo(imgBytes, args[0])

	// process
	src, _, _ := image.Decode(bytes.NewReader(imgBytes))
	newImg := imaging.FlipH(src)

	// output
	buffer := new(bytes.Buffer)
	imaging.Encode(buffer, newImg, imaging.JPEG)
	switch args[1].String() {
	case "string":
		return base64.StdEncoding.EncodeToString(buffer.Bytes())
	default:
		jsArray := js.Global().Get("Uint8Array").New(args[2].Get("length").Int())
		js.CopyBytesToJS(jsArray, buffer.Bytes())
		return jsArray
	}
}

func GrayImage(this js.Value, args []js.Value) any {
	// retrieve args
	imgBytesLen := args[0].Get("length").Int()
	imgBytes := make([]byte, imgBytesLen)
	js.CopyBytesToGo(imgBytes, args[0])

	// process
	src, _, _ := image.Decode(bytes.NewReader(imgBytes))
	newImg := imaging.Grayscale(src)

	// output
	buffer := new(bytes.Buffer)
	imaging.Encode(buffer, newImg, imaging.JPEG)
	switch args[1].String() {
	case "string":
		return base64.StdEncoding.EncodeToString(buffer.Bytes())
	default:
		jsArray := js.Global().Get("Uint8Array").New(args[2].Get("length").Int())
		js.CopyBytesToJS(jsArray, buffer.Bytes())
		return jsArray
	}
}

func KernelSmoothing(this js.Value, args []js.Value) any {
	// retrieve args
	rows := args[0].Int()
	cols := args[1].Int()
	data := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		data[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			data[i][j] = args[2].Index(i*cols + j).Float()
		}
	}

	// kernel smoothing
	kernel := [][]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	kernelSum := 16.0

	smoothed := make([][]float64, rows)
	for i := range smoothed {
		smoothed[i] = make([]float64, cols)
	}

	for i := 1; i < rows-1; i++ {
		for j := 1; j < cols-1; j++ {
			sum := 0.0
			for ki := -1; ki <= 1; ki++ {
				for kj := -1; kj <= 1; kj++ {
					sum += data[i+ki][j+kj] * kernel[ki+1][kj+1]
				}
			}
			smoothed[i][j] = sum / kernelSum
		}
	}

	// flatten the smoothed array to return to JS
	flattened := make([]float64, rows*cols)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			flattened[i*cols+j] = smoothed[i][j]
		}
	}

	jsArray := js.Global().Get("Float64Array").New(len(flattened))
	/*
		// convert flattened to []byte
		byteArray := make([]byte, len(flattened)*8)
		for i := range flattened {
			bits := math.Float64bits(flattened[i])
			binary.LittleEndian.PutUint64(byteArray[i*8:], bits)
		}
		js.CopyBytesToJS(jsArray, byteArray)
	*/
	// js.CopyBytesToJS(jsArray, float64SliceToByteSlice(flattened))
	// js.CopySliceToJS(jsArray, flattened)
	for i, v := range flattened {
		jsArray.SetIndex(i, v)
	}
	return jsArray
}

func KernelDensityEstimation(this js.Value, args []js.Value) any {
	start := time.Now()

	// retrieve args
	kdeSize := args[0].Int()
	dataLen := args[1].Get("length").Int()
	data := make([]int32, dataLen)
	for i := 0; i < dataLen; i++ {
		data[i] = int32(args[1].Index(i).Int())
	}

	// kernel density estimation
	xiData := make([]int32, kdeSize)
	for i := 0; i < kdeSize; i++ {
		xiData[i] = int32(i)
	}

	kdeData := make([]float64, kdeSize)
	N := float64(len(data))

	for i := 0; i < len(xiData); i++ {
		temp := 0.0
		for j := 0; j < len(data); j++ {
			temp += GaussKDE(xiData[i], data[j])
		}
		kdeData[i] = (1 / N) * temp
	}

	// create a Float64Array in JavaScript and copy the data
	jsArray := js.Global().Get("Float64Array").New(len(kdeData))
	for i, v := range kdeData {
		jsArray.SetIndex(i, v)
	}

	elapsed := time.Since(start)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	fmt.Printf("KernelDensityEstimation took %s and used %d bytes\n", elapsed, memStats.Alloc)

	return jsArray
}

func GaussKDE(xi, x int32) float64 {
	return (1 / math.Sqrt(2*math.Pi)) * math.Exp(math.Pow(float64(xi-x), 2)/-2)
}

/*
// Convert []float64 to []byte
func float64SliceToByteSlice(floats []float64) []byte {
	buf := make([]byte, len(floats)*8)
	for i, f := range floats {
		bits := math.Float64bits(f)
		binary.LittleEndian.PutUint64(buf[i*8:], bits)
	}
	return buf
}
*/

func main() {
	js.Global().Set("Resize", js.FuncOf(Resize))
	js.Global().Set("FlipImageH", js.FuncOf(FlipImageH))
	js.Global().Set("GrayImage", js.FuncOf(GrayImage))
	js.Global().Set("KernelSmoothing", js.FuncOf(KernelSmoothing))
	js.Global().Set("KernelDensityEstimation", js.FuncOf(KernelDensityEstimation))
	select {}
}
