import { writable } from "svelte/store";

export type UserToken = {
  token?: string,
  username?: string,
  name?: string,
  email?: string,
  permissions?: { name: string }[]
}

export const createSessionStore = (initialToken: UserToken) => writable(initialToken);
