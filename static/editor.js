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
const editorMain = document.getElementById("editor-main");
const editorPreview = document.getElementById("editor-preview");

outputDiv.innerHTML = go_parseMarkdown(inputDiv.value);

inputDiv.addEventListener("input", (e) => {
  e.preventDefault();

  outputDiv.innerHTML = go_parseMarkdown(inputDiv.value);
});

tabEditorLink.addEventListener("click", (e) => {
  e.preventDefault();
  editorMain.classList.remove("hide-touch");
  editorPreview.classList.add("hide-touch");
});

tabPreviewLink.addEventListener("click", (e) => {
  e.preventDefault();
  editorMain.classList.add("hide-touch");
  editorPreview.classList.remove("hide-touch");
});
