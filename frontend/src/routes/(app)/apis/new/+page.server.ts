import { error, redirect, type Cookies } from '@sveltejs/kit';
import { OK, ZodError, z } from 'zod';

async function callApi(url: string, access_token: string, body?: string) {
  return await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${access_token}`,
    },
    body: body
  });
}

async function callSaveApiEndpoint(validatedForm: {
  api_key_name?: string,
  api_key?: string,
  api_secret?: string,
  exchange?: string,
}, access_token: string) {
  let resp = await callApi(`http://${process.env.BACKEND_PRIVATE_HOST}/user/api_keys`, access_token, JSON.stringify(validatedForm));
  if (!resp.ok) {
    let data = await resp.json();
      return {
        detail: data.message,
        ok: false,
        values: validatedForm
      };
  }
  return {
    ok: true,
  }
}

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

    return await callSaveApiEndpoint({
      api_key_name: name,
      api_key: api_key,
      api_secret: api_secret,
      exchange: exchange,
    }, access_token);
  }
}
