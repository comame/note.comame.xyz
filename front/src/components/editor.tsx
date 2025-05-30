import { useState } from "react";
import { parse } from "../lib/markdown";
import "./editor.css";

interface props {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
}

export default function Editor({ onSubmit }: props) {
  const [text, setText] = useState("");
  const [markdown, setMarkdown] = useState("");

  const [tab, setTab] = useState<"editor" | "preview">("editor");

  return (
    <form
      id="editor-root"
      onSubmit={(e) => {
        stopPreventUnload();
        onSubmit(e);
      }}
      onChange={() => {
        startPreventUnload();
      }}
    >
      <nav id="editor-tab">
        <ul>
          <li>
            <a href="#" id="tab-editor" onClick={() => setTab("editor")}>
              Editor
            </a>
          </li>
          <li>
            <a href="#" id="tab-preview" onClick={() => setTab("preview")}>
              Preview
            </a>
          </li>
        </ul>
      </nav>

      <div id="editor-main" className={tab === "editor" ? "" : "hide-touch"}>
        <input id="title" name="title" placeholder="タイトル" required />
        <textarea
          required
          id="input"
          name="input"
          placeholder="本文"
          value={text}
          onChange={async (e) => {
            const md = e.currentTarget.value;
            setText(md);
            const parsed = await parse(md);
            setMarkdown(parsed);
          }}
        ></textarea>
      </div>
      <div
        id="editor-preview"
        className={tab === "preview" ? "" : "hide-touch"}
      >
        <div
          id="output"
          className="post-html"
          dangerouslySetInnerHTML={{ __html: markdown }}
        ></div>
      </div>

      <div id="control">
        <select name="visibility">
          <option value="0">非公開</option>
          <option value="1">限定公開</option>
          <option value="2">公開</option>
        </select>
        <button id="submit">SAVE</button>
      </div>
    </form>
  );
}

function beforeUnloadHandler(e: BeforeUnloadEvent) {
  e.preventDefault();
  e.returnValue = "";
}

function startPreventUnload() {
  window.addEventListener("beforeunload", beforeUnloadHandler);
}

function stopPreventUnload() {
  window.removeEventListener("beforeunload", beforeUnloadHandler);
}
