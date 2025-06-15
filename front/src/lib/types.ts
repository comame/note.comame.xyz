export type permission = "public" | "private" | "url";

export type goArr<T> = null | T[];

export type post = {
  id: number;
  url_key: string;
  createdDatetime: string;
  updatedDatetime: string;
  title: string;
  permissionInherited: boolean;
  resolvedPermission: permission;
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
  return `/posts/${post.url_key}`;
}
