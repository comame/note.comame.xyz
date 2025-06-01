import { useEffect, useState, type FormEvent } from "react";
import { parse } from "../lib/markdown";
import "./editor.css";
import type { permission, postConfig } from "../lib/types";
import BracketButton from "./bracket_button";
import Post from "./post";

interface props {
  onSubmit: (postConfig: postConfig) => void;
  editPost?: {
    title: string;
    text: string;
    permission: permission;
    id: number;
  };
  demoMode?: boolean;
}

export default function Editor({ onSubmit, editPost, demoMode }: props) {
  const draftID = getDraftID(editPost, demoMode);
  const draft = loadDraft(draftID);

  const [text, setText] = useState(draft?.text ?? editPost?.text ?? "");
  const [title, setTitle] = useState(draft?.title ?? editPost?.title ?? "");
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

  const onSubmitHandler = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    stopPreventUnload();

    if (demoMode) {
      return;
    }
    deleteDraft(draftID);
    onSubmit({
      title,
      text,
      permission: e.currentTarget.permission.value as permission,
      id: editPost?.id,
    });
  };

  return (
    <form
      id="editor-root"
      onSubmit={onSubmitHandler}
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
          value={title}
          onChange={(e) => {
            setTitle(e.currentTarget.value);
            saveDraft(e.currentTarget.value, text, draftID);
          }}
          required
        />
        <textarea
          required
          id="input"
          name="input"
          placeholder="本文"
          value={text}
          onChange={async (e) => {
            setText(e.currentTarget.value);
            saveDraft(title, e.currentTarget.value, draftID);
            const md = e.currentTarget.value;
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

function saveDraft(title: string, text: string, draftID: string | null) {
  if (!draftID) {
    return;
  }

  const key = `draft-${draftID}`;
  localStorage.setItem(key, JSON.stringify({ title, text }));
}

function loadDraft(
  draftID: string | null
): { title: string; text: string } | null {
  if (!draftID) {
    return null;
  }
  const key = `draft-${draftID}`;
  const draft = localStorage.getItem(key);
  if (draft) {
    return JSON.parse(draft);
  }
  return null;
}

function deleteDraft(draftID: string | null) {
  const key = `draft-${draftID}`;
  localStorage.removeItem(key);
}

function getDraftID(
  editPost?: { title: string; text: string; id: number } | undefined,
  demoMode?: boolean
): string | null {
  if (demoMode) {
    return null;
  }
  if (editPost) {
    return "" + editPost.id;
  } else {
    return "new";
  }
}
