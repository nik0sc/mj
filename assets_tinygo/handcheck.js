"use strict";

let wasm;
let eventCount = 0;
const go = new Go();

const readyInput = function() {
    if (eventCount < 1) {
        eventCount++;
    }

    const resultEl = document.getElementById("result");
    const errorEl = document.getElementById("error");
    const strategyEl = document.getElementById("selStrategy");
    const strategies = {
        "opt": optCheck,
        "optcnt": optCountCheck,
        "greedy": greedyCheck
    };

    document.getElementById("btnCheck").onclick = () => {
        const h = document.getElementById("inpHand").value;
        const split = document.getElementById("chkSplit").checked;
        const memo = document.getElementById("chkMemo").checked;
        const strategy = strategies[strategyEl.value];
        if (typeof strategy === undefined) {
            console.error("unknown strategy: " + strategyEl.value);
            return
        }

        strategy(h, (result, error) => {
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
