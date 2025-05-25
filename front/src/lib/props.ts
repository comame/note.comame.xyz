type globalProps = {
  Title: string;
  IsLoggedIn: boolean;
  Breadcrumbs: {
    Label: string;
    Location: string;
  }[];
  Page: pages;
  PageData: any;
};

type pages =
  | "TopPage"
  | "AllPostsPage"
  | "NewPostPage"
  | "EditPostPage"
  | "PostPage"
  | "NotFoundPage";

let globalPropsCache: globalProps | null = null;

export function getGlobalProps(): globalProps {
  if (globalPropsCache === null) {
    const element = document.querySelector("meta[name=template-props]");
    if (!element) {
      throw new Error("プロパティが渡されてない");
    }

    const json = element.getAttribute("content");
    if (!json) {
      throw new Error("プロパティが渡されてない");
    }

    globalPropsCache = JSON.parse(json);
  }

  return globalPropsCache!;
}
