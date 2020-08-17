document.addEventListener("DOMContentLoaded", function() {
    setup()
});

function setup() {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
        .then((result) => {
            go.run(result.instance);
        });
}