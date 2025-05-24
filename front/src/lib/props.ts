type globalProps = {
  Title: string;
  IsLoggedIn: boolean;
  OgDescription: string;
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
  | "PostPage";

let globalPropsCache: globalProps | null = null;

// TODO: test purpose
globalPropsCache = {
  Title: "dummy title",
  IsLoggedIn: true,
  OgDescription: "test",
  Page: "TopPage",
  Breadcrumbs: [
    {
      Label: "AAA",
      Location: "#aaa",
    },
    {
      Label: "BBB",
      Location: "#bbb",
    },
    {
      Label: "CCC",
      Location: "#ccc",
    },
  ],
  PageData: {},
};

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
