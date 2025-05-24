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
  const { Page: currentPage, PageData, IsLoggedIn } = getGlobalProps();

  switch (currentPage) {
    case "TopPage":
      return <TopPage pageData={PageData} />;
    case "AllPostsPage":
      return <AllPostsPage pageData={PageData} isLoggedIn={IsLoggedIn} />;
    case "PostPage":
      return <PostPage pageData={PageData} />;
    case "NewPostPage":
      return <NewPostPage pageData={PageData} />;
    case "EditPostPage":
      return <EditPostPage pageData={PageData} />;
  }
}

function App() {
  const { Breadcrumbs, IsLoggedIn } = getGlobalProps();

  const breadcrumbs = Breadcrumbs.map((v) => ({
    location: v.Location,
    label: v.Label,
  }));

  return (
    <div>
      <Header isLoggedIn={IsLoggedIn} breadcrumbs={breadcrumbs} />
      <Page />
    </div>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
