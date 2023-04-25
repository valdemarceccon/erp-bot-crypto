import { error, redirect } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const POST: RequestHandler = async ({cookies, request, fetch}) => {
  let fd = await request.formData();
  let username = fd.get("username");
  let password = fd.get("password");
  if (!username || !password) {
    throw error(401, {message: "username and password is mandatory"});
  }

  let resp = await fetch('http://backend:8000/auth/token', {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      "username": username.toString(),
      "password": password.toString()
    })
  });


  if (!resp.ok) {
    let data = await resp.json();
    throw error(resp.status, data.message)
  }

  let token_resp = await resp.json();
  cookies.set("access_token", token_resp.access_token, {
    path: "/"
  })

  throw redirect(301, "/");
}
