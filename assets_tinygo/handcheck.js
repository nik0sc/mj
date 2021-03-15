"use strict";

let wasm;
const go = new Go();

const readyInput = function() {
    const resultEl = document.getElementById("result");
    const errorEl = document.getElementById("error");
    const strategyEl = document.getElementById("selStrategy");
    const inpHandEl = document.getElementById("inpHand");
    const btnCheckEl = document.getElementById("btnCheck");
    btnCheckEl.onclick = () => {
        const h = inpHandEl.value;
        const split = document.getElementById("chkSplit").checked;
        const memo = document.getElementById("chkMemo").checked;
        const strategy = strategyEl.value;

        checkHand(strategy, h, (result, error) => {
            if (error != null) {
                resultEl.textContent = "";
                errorEl.textContent = error;
                console.error(error);
            } else {
                resultEl.textContent = result;
                errorEl.textContent = "";
            }
        }, split, memo);
    };
    document.getElementById("btnCheck").disabled = false;
    errorEl.textContent = "";

    document.querySelectorAll(".tryout").forEach((el) => {
        el.addEventListener("click", () => {
            inpHandEl.value = el.textContent;
            btnCheckEl.click();
        })
    })

    console.log("input ready");
}

const wasmCallback = function (obj) {
    wasm = obj.instance;
    go.run(wasm);
    console.log("wasm ready");
    readyInput();
};

if ('instantiateStreaming' in WebAssembly) {
    WebAssembly.instantiateStreaming(fetch("handcheck.wasm"), go.importObject).then(wasmCallback)
} else {
    fetch("handcheck.wasm").then(
        resp => resp.arrayBuffer()
    ).then(
        bytes => WebAssembly.instantiate(bytes, go.importObject).then(wasmCallback)
    );
}

window.addEventListener("DOMContentLoaded", readyInput);
