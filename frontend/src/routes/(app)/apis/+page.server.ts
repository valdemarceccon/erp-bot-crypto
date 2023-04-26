import { redirect } from '@sveltejs/kit';

export type ApiListResp = {
  id: number,
  name: string,
  api_key: string,
  exchange: string,
  status: number
};

export async function load({cookies, fetch}) {
  const token = cookies.get("access_token");

  if (!token) {
    redirect(301, "/");
  }

  let api_keys_resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/users/api_keys/`, {
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

  let res_api_list: ApiListResp[] = await api_keys_resp.json();

  return {
    success: true,
    api_list: res_api_list
  }

}
