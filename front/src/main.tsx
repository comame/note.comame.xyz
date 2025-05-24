import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import AllPostsPage from "./pages/AllPostsPage";
import EditPostPage from "./pages/EditPostPage";
import NewPostPage from "./pages/NewPostPage";
import PostPage from "./pages/PostPage";
import TopPage from "./pages/TopPage";
import Header from "./components/header";
import { getGlobalProps } from "./lib/props";

function Page() {
  const { Page: currentPage } = getGlobalProps();

  switch (currentPage) {
    case "TopPage":
      return <TopPage />;
    case "AllPostsPage":
      return <AllPostsPage />;
    case "PostPage":
      return <PostPage />;
    case "NewPostPage":
      return <NewPostPage />;
    case "EditPostPage":
      return <EditPostPage />;
  }
}

function App() {
  return (
    <div>
      <Header />
      <Page />
    </div>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
