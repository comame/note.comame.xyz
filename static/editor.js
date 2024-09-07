import "/out/dist/wasm_exec.js";

const go = new Go();
window._go = go;

const { instance } = await WebAssembly.instantiateStreaming(
  fetch("/out/dist/goapp.wasm"),
  go.importObject
);

go.run(instance);

const inputDiv = document.getElementById("input");
const outputDiv = document.getElementById("output");
const tabEditorLink = document.getElementById("tab-editor");
const tabPreviewLink = document.getElementById("tab-preview");

outputDiv.innerHTML = go_parseMarkdown(inputDiv.value);

inputDiv.addEventListener("input", (e) => {
  e.preventDefault();

  outputDiv.innerHTML = go_parseMarkdown(inputDiv.value);
});

tabEditorLink.addEventListener("click", () => {
  inputDiv.classList.remove("hide-touch");
  outputDiv.classList.add("hide-touch");
});

tabPreviewLink.addEventListener("click", () => {
  inputDiv.classList.add("hide-touch");
  outputDiv.classList.remove("hide-touch");
});
