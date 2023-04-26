import { error, redirect } from '@sveltejs/kit';
import { OK, ZodError, z } from 'zod';

async function callSaveApiEndpoint(validatedForm: {
  name?: string,
  api_key?: string,
  api_secret?: string,
  exchange?: string,
}, access_token: string) {
  let resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/users/api_key`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${access_token}`,
    },
    body: JSON.stringify(validatedForm)
  });
  if (!resp.ok) {
    let data = await resp.json();
    if (resp.status >= 400 && resp.status < 500) {
      return {
        detail: data.detail,
        ok: false,
        values: validatedForm
      };
    }
    throw error(resp.status, data.message)
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
      name: name,
      api_key: api_key,
      api_secret: api_secret,
      exchange: exchange,
    }, access_token);
  }
}
