import { localStorageStore } from "@skeletonlabs/skeleton"
import type { Writable } from "svelte/store"

export type UserToken = {
  token?: string,
  username?: string,
  name?: string,
  email?: string,
  permissions?: { name: string }[]
}

export const userTokenStore: Writable<UserToken> = localStorageStore('user', {});
