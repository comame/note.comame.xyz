import "/out/dist/wasm_exec.js";

const go = new Go();

const { instance } = await WebAssembly.instantiateStreaming(
  fetch("/out/dist/goapp.wasm"),
  go.importObject
);

go.run(instance);

const inputDiv = document.getElementById("input");
const outputDiv = document.getElementById("output");
const titleInput = document.getElementById("title");
const tabEditorLink = document.getElementById("tab-editor");
const tabPreviewLink = document.getElementById("tab-preview");
const editorMain = document.getElementById("editor-main");
const editorPreview = document.getElementById("editor-preview");
const form = document.getElementById("editor-root");
const isDemoMeta = document.querySelector("meta[name=is-demo]");

const draft = getDraftForCurrentPage();
if (draft !== null && window.confirm("下書きを読み込みますか？")) {
  titleInput.value = draft.title;
  inputDiv.value = draft.text;
}

outputDiv.innerHTML = go_parseMarkdown(inputDiv.value);

inputDiv.addEventListener("input", (e) => {
  e.preventDefault();

  const fd = new FormData(form);
  saveDraftForCurrentPage(fd.get("title"), fd.get("input"));
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
  clearDraftForCurrentPage();
  location.replace(js["location"]);
});

function isDemo() {
  return isDemoMeta.getAttribute("value") === "true";
}

function getEditingPostID() {
  const fd = new FormData(form);
  const id = Number.parseInt(fd.get("id"), 10);

  if (id === 0) {
    return null;
  }
  return id;
}

function getLocalStorageKey() {
  const id = getEditingPostID();

  if (id === null) {
    return "draft:new";
  }
  return "draft:" + id;
}

/**
 * @param {string} text
 * @param {string} title
 */
function saveDraftForCurrentPage(title, text) {
  if (isDemo()) {
    return;
  }

  const key = getLocalStorageKey();

  localStorage.setItem(key, JSON.stringify({ title, text }));
}

function clearDraftForCurrentPage() {
  const key = getLocalStorageKey();
  localStorage.removeItem(key);
}

function getDraftForCurrentPage() {
  if (isDemo()) {
    return null;
  }

  const key = getLocalStorageKey();

  const draft = localStorage.getItem(key);
  if (draft === null) {
    return null;
  }

  try {
    const data = JSON.parse(draft);
    return {
      title: data["title"],
      text: data["text"],
    };
  } catch (_) {
    return null;
  }
}
