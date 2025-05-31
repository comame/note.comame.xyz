import BracketButton from "../../components/bracket_button";
import Post from "../../components/post";
import StatusBadge from "../../components/status_badge";
import type { post } from "../../lib/types";
import "./index.css";

type pageData = {
  post: post;
};

export default function PostPage({
  pageData,
  isLoggedIn,
}: {
  pageData: pageData;
  isLoggedIn: boolean;
}) {
  const { post } = pageData;

  const editLink = `/edit/post/${post.id}`;

  return (
    <div className="pages-post">
      <div id="post">
        <div className="metadata-title">
          <h1 className="title">{post.title}</h1>
          <span className="c-visibility">
            <StatusBadge size="m" status={post.permission} />
          </span>
          {isLoggedIn && (
            <span className="edit-link">
              <BracketButton values={[{ label: "編集", to: editLink }]} />
            </span>
          )}
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
        <Post html={post.html} />
      </div>
    </div>
  );
}
