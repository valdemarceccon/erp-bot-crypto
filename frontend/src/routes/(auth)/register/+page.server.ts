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
      let data = await resp.json();
      if (resp.status >= 400 && resp.status < 500) {
        return {
          message: data.message,
          ok: false,
          values: {
            username: username,
            password: password,
            name: name,
            email: email,
            password_confirm: password_confirm,
          }
        };
      }
      throw error(resp.status, data)
    }
    let token_resp = await resp.json();
    cookies.set("access_token", token_resp.token, {
      path: "/"
    })

    throw redirect(301, "/login");
  }
}
