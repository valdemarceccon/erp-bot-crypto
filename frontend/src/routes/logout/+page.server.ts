import { redirect } from '@sveltejs/kit';

export const actions = {
  default: ({ cookies }) => {
    cookies.delete("access_token", {
      path: "/"
    });

    throw redirect(301, "/");
  }
}
