import { useEffect, useState } from "react";
import { parse } from "../lib/markdown";
import "./editor.css";
import type { permission } from "../lib/types";

interface props {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  editPost?: {
    title: string;
    text: string;
    permission: permission;
    id: string;
  };
  demoMode?: boolean;
}

export default function Editor({ onSubmit, editPost, demoMode }: props) {
  const [text, setText] = useState("");
  const [markdown, setMarkdown] = useState("");

  const [tab, setTab] = useState<"editor" | "preview">("editor");

  useEffect(() => {
    if (!editPost) {
      return;
    }
    parse(editPost.text).then((md) => {
      setMarkdown(md);
    });
  }, [editPost]);

  let defaultPermission = undefined;
  if (editPost) {
    defaultPermission = permissionToVisibility(editPost.permission);
  }

  return (
    <form
      id="editor-root"
      onSubmit={(e) => {
        stopPreventUnload();

        if (demoMode) {
          return;
        }
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
        <input
          id="title"
          name="title"
          placeholder="タイトル"
          defaultValue={editPost?.title ?? undefined}
          required
        />
        <textarea
          required
          id="input"
          name="input"
          placeholder="本文"
          value={text}
          defaultValue={editPost?.text ?? undefined}
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

      {!demoMode && (
        <div id="control">
          <select name="visibility" defaultValue={defaultPermission}>
            <option value="0">非公開</option>
            <option value="1">限定公開</option>
            <option value="2">公開</option>
          </select>
          <button id="submit">SAVE</button>
        </div>
      )}

      {editPost && <input type="hidden" name="id" value={editPost.id} />}
    </form>
  );
}

function permissionToVisibility(p: permission): number {
  switch (p) {
    case "private":
      return 0;
    case "url":
      return 1;
    case "public":
      return 2;
  }
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
