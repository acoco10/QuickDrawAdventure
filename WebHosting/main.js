fetch('QDA.wasm')
    .then(response => response.arrayBuffer())
    .then(bytes => WebAssembly.instantiate(bytes))
    .then(result => {
        document.getElementById('output').textContent = "WASM Loaded!";
        console.log(result.instance);
    })
    .catch(console.error);