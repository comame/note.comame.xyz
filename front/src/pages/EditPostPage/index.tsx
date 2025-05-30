import type { FormEvent } from "react";
import Editor from "../../components/editor";
import type { post } from "../../lib/types";
import "./index.css";

type pageData = {
  post: post;
};

export default function EditPostPage({ pageData }: { pageData: pageData }) {
  return (
    <div className="pages-edit-post">
      <Editor onSubmit={onSubmit} editPost={pageData.post} />
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

  const res = await fetch(`/edit/post/${json.id}`, {
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
