import "./index.css";
import Editor from "../../components/editor";
import specMD from "../../../../internal/md/spec.md?raw";

export default function NewPostPage() {
  const post = {
    text: specMD,
    title: "Markdown エディタ",
    resolvedPermission: "public",
    id: 0,
  } as const;

  return (
    <div className="pages-demo-editor">
      <Editor onSubmit={() => {}} demoMode editPost={post} />
    </div>
  );
}
