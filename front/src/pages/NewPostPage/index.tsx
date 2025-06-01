import "./index.css";
import Editor from "../../components/editor";
import type { postConfig } from "../../lib/types";

type pageData = {};

export default function NewPostPage({ pageData: _ }: { pageData: pageData }) {
  return (
    <div className="pages-new-post">
      <Editor onSubmit={onSubmit} />
    </div>
  );
}

async function onSubmit(postConfig: postConfig) {
  const res = await fetch("/post/create", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(postConfig),
    redirect: "manual",
  });

  if (!res.ok) {
    return;
  }

  const js = await res.json();
  location.replace(js["location"]);
}
