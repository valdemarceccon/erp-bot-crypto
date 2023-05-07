import { fail, redirect } from '@sveltejs/kit';

export type ApiListAdminResp = {
  name: string;
  api_key: string;
  user_id: number;
  exchange: string;
  status: number;
  id: number;
  username: string;
  api_secret: string;
};

export const actions = {
  toggleStatus: async ({ cookies, fetch, url, request }) => {
    const token = cookies.get("access_token");

    if (!token) {
      redirect(301, "/");
    }
    let fd = await request.formData();
    let id = fd.get("id")?.toString();
    let client_id = fd.get("client_id")?.toString();

    if (!id) {
      fail(403, {
        success: false,
        message: "invalid request"
      })
    }
    let api_keys_resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/api_keys/admin-toggle/${client_id}/${id}`, {
      method: "PATCH",
      headers: {
        "Authorization": `Bearer ${token}`,
        "Content-Type": "application/json"
      }
    });

    if (!api_keys_resp.ok) {
      return await api_keys_resp.json();
    }

    let res_api_list: ApiListAdminResp[] = await api_keys_resp.json();

    return {
      success: true,
      api_list: res_api_list
    }
  }
}

export async function load({ cookies, fetch }) {
  const token = cookies.get("access_token");

  if (!token) {
    redirect(301, "/");
  }

  let api_keys_resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/api_keys/all`, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
    }
  });

  if (!api_keys_resp.ok) {
    let a = await api_keys_resp.json();
    if (api_keys_resp.status < 500) {
      return {
        success: false,
        message: a.detail
      }
    }
  }

  let res_api_list: ApiListAdminResp[] = await api_keys_resp.json();

  return {
    success: true,
    api_list: res_api_list
  }

}
