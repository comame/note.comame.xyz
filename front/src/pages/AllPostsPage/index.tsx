import "./index.css";
import PostTable from "../../components/post_table";
import type { post } from "../../lib/types";

type pageData = {
  posts: post[];
};

type props = {
  pageData: pageData;
  isLoggedIn: boolean;
};

export default function AllPostsPage({ pageData, isLoggedIn }: props) {
  const { posts } = pageData;

  return (
    <div className="page-all-posts">
      <PostTable posts={posts} isLoggedIn={isLoggedIn} />
    </div>
  );
}
