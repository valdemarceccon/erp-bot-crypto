import { redirect } from '@sveltejs/kit';
import './$types';

export function load({ cookies }) {
  cookies.delete("access_token", {
    path: "/"
  });

  throw redirect(301, "/");
}
