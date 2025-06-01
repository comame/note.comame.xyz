import Editor from "../../components/editor";
import type { post, postConfig } from "../../lib/types";
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

async function onSubmit(postConfig: postConfig) {
  const res = await fetch(`/edit/post/${postConfig.id}`, {
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
