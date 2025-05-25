export type permission = "public" | "private" | "url";

export type post = {
  id: string;
  urlKey: string;
  createdDatetime: string;
  updatedDatetime: string;
  title: string;
  permissionInherited: boolean;
  permission: permission;
  text: string;
  html: string;
};
