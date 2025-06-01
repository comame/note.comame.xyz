export type permission = "public" | "private" | "url";

export type post = {
  id: number;
  url_key: string;
  createdDatetime: string;
  updatedDatetime: string;
  title: string;
  permissionInherited: boolean;
  permission: permission;
  text: string;
  html: string;
};

export type postConfig = {
  id?: number;
  text: string;
  title: string;
  permission: permission;
};

export function getPostURL(post: post): string {
  switch (post.permission) {
    case "public":
      return `/posts/public/${post.url_key}`;
    case "private":
      return `/posts/private/${post.url_key}`;
    case "url":
      return `/posts/unlisted/${post.url_key}`;
  }
}
