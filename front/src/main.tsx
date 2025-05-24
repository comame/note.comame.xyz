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
import PostTable from "./components/post_table";

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
  const { Breadcrumbs, IsLoggedIn } = getGlobalProps();

  const breadcrumbs = Breadcrumbs.map((v) => ({
    location: v.Location,
    label: v.Label,
  }));

  return (
    <div>
      <Header isLoggedIn={IsLoggedIn} breadcrumbs={breadcrumbs} />
      <Page />

      <PostTable
        posts={[
          {
            urlKey: "foo",
            createdDatetime: "2024-01-01 00:00:00",
            updatedDatetime: "2025-01-01 00:00:00",
            title: "タイトルだよ",
            permission: "public",
            permissionInherited: false,
            text: "",
            html: "",
          },
          {
            urlKey: "foobar",
            createdDatetime: "2024-01-01 00:00:00",
            updatedDatetime: "2025-12-01 00:00:00",
            title: "タイトルだよaaaaaaaaaaa",
            permission: "public",
            permissionInherited: false,
            text: "",
            html: "",
          },
        ]}
        isLoggedIn={IsLoggedIn}
      />
    </div>
  );
}

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
