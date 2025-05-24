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
  const { Page: currentPage, PageData } = getGlobalProps();

  switch (currentPage) {
    case "TopPage":
      return <TopPage pageData={PageData} />;
    case "AllPostsPage":
      return <AllPostsPage pageData={PageData} />;
    case "PostPage":
      return <PostPage pageData={PageData} />;
    case "NewPostPage":
      return <NewPostPage pageData={PageData} />;
    case "EditPostPage":
      return <EditPostPage pageData={PageData} />;
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
