import { getPostURL, type post } from "../lib/types";
import BracketButton from "./bracket_button";
import StatusBadge from "./status_badge";
import "./post_table.css";

type props = {
  posts: post[];
  isLoggedIn: boolean;
};

export default function PostTable({ posts, isLoggedIn }: props) {
  return (
    <table
      className={`components-post-table ` + (isLoggedIn ? "is-logged-in" : "")}
    >
      <thead>
        <tr>
          <div>
            <td className="title">タイトル</td>
            <td>更新</td>
          </div>
          <div>
            {isLoggedIn && <td></td>}
            {isLoggedIn && <td></td>}
          </div>
        </tr>
      </thead>
      <tbody>
        {posts.map((post) => (
          <tr key={getPostURL(post)}>
            <div>
              <td className="title">
                <a href={getPostURL(post)}>{post.title}</a>
              </td>
              <td>{formatDate(post.updatedDatetime)}</td>
            </div>
            <div>
              {isLoggedIn && (
                <td>
                  <StatusBadge
                    size="m"
                    status={post.permission}
                    inherit={post.permissionInherited}
                  />
                </td>
              )}
              {isLoggedIn && (
                <td>
                  <BracketButton
                    values={[
                      {
                        label: "編集",
                        to: "#",
                      },
                      {
                        label: "削除",
                        onClick: () => {
                          deletePost(post.id);
                        },
                      },
                    ]}
                  />
                </td>
              )}
            </div>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function formatDate(datetime: string) {
  const d = new Date(datetime);

  const now = Date.now();
  const diff = now - d.getTime();

  if (diff >= 1000 * 86400 * 330) {
    return `${d.getFullYear}年`;
  }

  return `${d.getMonth() + 1}月${d.getDate()}日`;
}

function deletePost(postID: string) {
  if (!confirm("削除しますか")) {
    return;
  }

  fetch("/delete/post/" + postID, {
    method: "POST",
  }).then((res) => {
    if (!res.ok) {
      return;
    }

    location.reload();
  });
}
