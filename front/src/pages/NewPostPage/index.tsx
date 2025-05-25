import "./index.css";
import { parse } from "../../lib/markdown";
import { useEffect, useState, type FormEvent } from "react";

type pageData = {};

export default function NewPostPage({ pageData: _ }: { pageData: pageData }) {
  const [text, setText] = useState("");
  const [markdown, setMarkdown] = useState("");

  const [tab, setTab] = useState<"editor" | "preview">("editor");

  useEffect(() => {
    const handler = (e: BeforeUnloadEvent) => {
      e.preventDefault();
      e.returnValue = "";
    };
    window.addEventListener("beforeunload", handler);
    return () => {
      window.removeEventListener("beforeunload", handler);
    };
  }, []);

  return (
    <div className="pages-new-post">
      <form id="editor-root" onSubmit={onSubmit}>
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
    </div>
  );
}

async function onSubmit(e: FormEvent<HTMLFormElement>) {
  e.preventDefault();

  const form = e.currentTarget as HTMLFormElement;

  const fd = new FormData(form);
  const json = {
    visibility: Number.parseInt(fd.get("visibility") as string, 10),
    text: fd.get("input"),
    title: fd.get("title"),
    id: Number.parseInt(fd.get("id") as string, 10),
  };

  const res = await fetch("/post/create", {
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
}
