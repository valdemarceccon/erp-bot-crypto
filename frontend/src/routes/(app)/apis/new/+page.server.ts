import { error, redirect, type Cookies, json } from '@sveltejs/kit';
import { OK, ZodError, z } from 'zod';

export const actions = {
  default: async ({ cookies, request, fetch }) => {
    const access_token = cookies.get("access_token",);

    if (!access_token) {
      throw redirect(301, "/login")
    }

    let fd = await request.formData();
    let name = fd.get("name")?.toString();
    let api_key = fd.get("api_key")?.toString();
    let api_secret = fd.get("api_secret")?.toString();
    let exchange = fd.get("exchange")?.toString();

    let ret = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/api_keys`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${access_token}`,
      },
      body: JSON.stringify({
        api_key_name: name,
        api_key: api_key,
        api_secret: api_secret,
        exchange: exchange,
      })
    });

    if (!ret.ok)
      return (await ret.json());

    throw redirect(302, "/apis")
  }
}
