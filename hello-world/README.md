# Go WebAssembly README

This project demonstrates how to use Go to compile WebAssembly (Wasm) modules.

## Prerequisites

- Go 1.11 or later (current version: 1.24.1)
- A modern web browser

## Getting Started

1. **Install Go**: Follow the instructions on the [official Go website](https://golang.org/doc/install) to install Go.

2. **Enable WebAssembly support**: Ensure you have Go modules enabled by setting the `GO111MODULE` environment variable to `on`.

3. **Create a new project**:

    ```sh
    mkdir go-wasm
    cd go-wasm
    go mod init example.com/go-wasm
    ```

4. **Check Env And Install Package**:

    ```sh
    ls $(go env GOROOT)/lib/wasm/
    # go_js_wasm_exec         go_wasip1_wasm_exec     wasm_exec.js            wasm_exec_node.js
    cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .
    npm install -g simplehttpserver
    ```

5. **Write Go code**: Create a `main.go` file with the following content:

    ```go
    package main

    import (
        "syscall/js"
    )

    func main() {
        js.Global().Set("sayHello", js.FuncOf(sayHello))
        select {}
    }

    func sayHello(this js.Value, p []js.Value) interface{} {
        return "Hello, WebAssembly!"
    }
    ```

6. **Compile to WebAssembly**:

    ```sh
    GOOS=js GOARCH=wasm go build -o main.wasm
    ```

7. **Serve the WebAssembly module**: Create an `index.html` file to load the Wasm module:

    ```html
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>Go WebAssembly</title>
    </head>
    <body>
        <h1>Go WebAssembly Example</h1>
        <button onclick="sayHello()">Say Hello</button>
        <script>
            const go = new Go();
            WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
                go.run(result.instance);
            });
        </script>
    </body>
    </html>
    ```

8. **Run a local server**:

    ```sh
    go run $(go env GOROOT)/lib/wasm/serve
    ```

9. **Open your browser**: Navigate to `http://localhost:8080` to see the WebAssembly module in action.

## Resources

- [Go WebAssembly Documentation](https://golang.org/doc/go1.11#wasm)
- [WebAssembly Official Site](https://webassembly.org/)

## License

This project is licensed under the MIT License.
