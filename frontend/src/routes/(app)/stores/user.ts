import { writable, type Writable } from "svelte/store";

type UserInfo = {
  username?: string,
  name?: string,
  email?: string,
  access_token?: string
}

let a: UserInfo = {}

export const userStore: Writable<UserInfo> = writable(a);
