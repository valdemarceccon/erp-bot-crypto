import { type Actions, error, json, redirect } from "@sveltejs/kit";
export function load({cookies}) {
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
    if (!username || !password) {
      throw error(401, { message: "username and password is mandatory" });
    }

    let resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/auth/token`, {
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
      if (resp.status == 401) {
        return {
          detail: data.detail,
          ok: false,
          values: {
            username: username,
            password: password
          }
        };
      }
      throw error(resp.status, data.message)
    }

    let token_resp = await resp.json();
    cookies.set("access_token", token_resp.access_token, {
      path: "/"
    })

    throw redirect(301, "/");
  }
}
