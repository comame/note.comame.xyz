import StatusBadge from "../../components/status_badge";
import type { post } from "../../lib/types";
import "./index.css";

type pageData = {
  post: post;
};

export default function PostPage({ pageData }: { pageData: pageData }) {
  const { post } = pageData;

  console.log(post);

  return (
    <div className="pages-post">
      <div id="post">
        <div className="metadata-title">
          <h1 className="title">{post.title}</h1>
          <span className="c-visibility">
            <StatusBadge size="m" status={post.permission} />
          </span>
        </div>
        <ul className="time">
          <li>
            created:&nbsp;
            <time>{post.createdDatetime}</time>
          </li>
          <li>
            updated:&nbsp;
            <time>{post.updatedDatetime}</time>
          </li>
        </ul>
        <div
          className="post-html"
          dangerouslySetInnerHTML={{ __html: post.html }}
        ></div>
      </div>
    </div>
  );
}
