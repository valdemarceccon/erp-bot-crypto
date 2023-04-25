import { error, json, redirect, text,  } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({cookies}) => {
  cookies.delete("access_token", {
    path: "/"
  });

  throw redirect(301, "/");
}
