import { type Actions, error, json, redirect } from "@sveltejs/kit";
export function load({ cookies }) {
  let c = cookies.get("access_token");
  if (c) {
    throw redirect(301, "/");
  }
}

export const actions: Actions = {
  default: async ({ cookies, request, fetch }) => {
    let fd = await request.formData();

    let username = fd.get("username");
    let password = fd.get("password");
    let name = fd.get("name");
    let email = fd.get("email");
    let password_confirm = fd.get("password_confirm");

    let validation: { username?: string, password?: string, name?: string, email?: string } = {}

    if (!username) validation.username = "username is required";
    if (!password) validation.password = "password is required";
    if (!name) validation.name = "name is required";
    if (!email) validation.email = "email is required";
    if (password && password_confirm && password !== password_confirm) validation.password = "password confirmation must match"

    if (Object.keys(validation).length > 0) {
      return {
        ok: false,
        validation: validation
      }
    }

    let resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        "username": username!.toString(),
        "password": password!.toString(),
        "email": email!.toString(),
        "fullname": name!.toString(),
      })
    });

    if (!resp.ok) {
      return await resp.json();
    }
    let token_resp = await resp.json();
    cookies.set("access_token", token_resp.token, {
      path: "/"
    })

    throw redirect(301, "/login");
  }
}
