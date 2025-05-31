import "./index.css";
import { type FormEvent } from "react";
import Editor from "../../components/editor";
import type { permission, post } from "../../lib/types";

type pageData = {};

export default function NewPostPage({ pageData: _ }: { pageData: pageData }) {
  return (
    <div className="pages-new-post">
      <Editor onSubmit={onSubmit} />
    </div>
  );
}

async function onSubmit(e: FormEvent<HTMLFormElement>) {
  e.preventDefault();

  const form = e.currentTarget as HTMLFormElement;

  const fd = new FormData(form);
  const json: Partial<post> = {
    permission: fd.get("permission") as permission,
    text: fd.get("input") as string,
    title: fd.get("title") as string,
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
