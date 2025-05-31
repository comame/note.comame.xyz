import { useEffect, useState } from "react";
import { parse } from "../lib/markdown";
import "./editor.css";
import type { permission } from "../lib/types";
import BracketButton from "./bracket_button";
import Post from "./post";

interface props {
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
  editPost?: {
    title: string;
    text: string;
    permission: permission;
    id: number;
  };
  demoMode?: boolean;
}

export default function Editor({ onSubmit, editPost, demoMode }: props) {
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

  const defaultPermission = editPost?.permission ?? "private";

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
        <BracketButton
          values={[
            { label: "Editor", onClick: () => setTab("editor") },
            { label: "Preview", onClick: () => setTab("preview") },
          ]}
        />
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
          // value={text}
          defaultValue={editPost?.text ?? undefined}
          onChange={async (e) => {
            const md = e.currentTarget.value;
            // setText(md);
            const parsed = await parse(md);
            setMarkdown(parsed);
          }}
        ></textarea>
      </div>
      <div
        id="editor-preview"
        className={tab === "preview" ? "" : "hide-touch"}
      >
        <Post html={markdown} />
      </div>

      {!demoMode && (
        <div id="control">
          <select name="permission" defaultValue={defaultPermission}>
            <option value="private">非公開</option>
            <option value="url">限定公開</option>
            <option value="public">公開</option>
          </select>
          <button id="submit">SAVE</button>
        </div>
      )}

      {editPost && <input type="hidden" name="id" value={editPost.id} />}
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
