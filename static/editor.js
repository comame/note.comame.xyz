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
const form = document.getElementById("editor-root");

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

form.addEventListener("submit", async (e) => {
  e.preventDefault();
  const fd = new FormData(form);
  const json = {
    visibility: Number.parseInt(fd.get("visibility"), 10),
    text: fd.get("input"),
    title: fd.get("title"),
    id: Number.parseInt(fd.get("id"), 10),
    url_key: fd.get("url-key"),
  };

  const res = await fetch(form.target, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(json),
    redirect: "manual",
  });

  if (!res.ok) {
    return;
  }

  const js = await res.json();
  location.replace(js["location"]);
});
