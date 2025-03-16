self.importScripts('wasm_exec.js');

const go = new Go();

/*
const WASM_URL = 'main.wasm';

var wasm;
async (e) => {
if (!self.wasmInstance) {
    // const result = await WebAssembly.instantiateStreaming(fetch('asset/main.wasm'), go.importObject);
    /*
    const response = await fetch('tiny.wasm');
    if (!response.ok) {
        throw new Error(`Failed to fetch WebAssembly module. HTTP status code: ${response.status}`);
    }
    const bytes = await response.arrayBuffer();
    const result = await WebAssembly.instantiate(bytes, go.importObject);
    go.run(result.instance);
    self.wasmInstance = result.instance;
    //
    if ('instantiateStreaming' in WebAssembly) {
        WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
            wasm = obj.instance;
            go.run(wasm);
        })
    } else {
        fetch(WASM_URL).then(resp =>
            resp.arrayBuffer()
        ).then(bytes =>
            WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
                wasm = obj.instance;
                go.run(wasm);
            })
        )
    }  
}
}
*/



self.onmessage = async (e) => {
    const { type, width, height, srcImgBytes } = e.data;

    if (!self.wasmInstance) {
        // const result = await WebAssembly.instantiateStreaming(fetch('asset/main.wasm'), go.importObject);
        /*const headers = new Headers();
        headers.append('Content-Encoding', 'gzip');
        headers.append('Content-Type', 'application/gzip');
        const response = await fetch('main.wasm.gz', { headers: headers });
        */
        const response = await fetch('main.wasm');
        if (!response.ok) {
            throw new Error(`Failed to fetch WebAssembly module. HTTP status code: ${response.status}`);
        }
        const bytes = await response.arrayBuffer();
        const result = await WebAssembly.instantiate(bytes, go.importObject);
        go.run(result.instance);
        self.wasmInstance = result.instance;
    }

    let response;
    switch (type) {
        case 'resize:arraybuffer':
            // Go return Uint8Array, so we need to convert it to base64 string.
            const newImgBytes = Resize(width, height, srcImgBytes, 'arraybuffer');
            const base64str = btoa(new Uint8Array(newImgBytes).reduce((data, byte) => data + String.fromCharCode(byte), ''));
            response = { eventType: 'resize:arraybuffer', base64: base64str};
            break;
        case 'resize:string':
            // Go return base64 string, so we can directly use it in JS runtime.
            const newResizeBase64 = Resize(width, height, srcImgBytes, 'string');
            response = { eventType: 'resize:string', base64: newResizeBase64 };
            break;
        case 'flipH':
            const newFlipImageH = FlipImageH(srcImgBytes, 'string');
            response = { eventType: 'flipH', base64: newFlipImageH };
            break;
        case 'gray':
            const newGrayImage = GrayImage(srcImgBytes, 'string');
            response = { eventType: 'gray', base64: newGrayImage };
            break;
    }

    self.postMessage(response);
};
