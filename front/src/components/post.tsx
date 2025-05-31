import "./post.css";

interface props {
  html: string;
}

export default function Post({ html }: props) {
  return (
    <div
      className="components-post"
      dangerouslySetInnerHTML={{ __html: html }}
    ></div>
  );
}
