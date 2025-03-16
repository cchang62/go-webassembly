package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"syscall/js"

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

func main() {
	js.Global().Set("Resize", js.FuncOf(Resize))
	js.Global().Set("FlipImageH", js.FuncOf(FlipImageH))
	js.Global().Set("GrayImage", js.FuncOf(GrayImage))
	select {}
}
