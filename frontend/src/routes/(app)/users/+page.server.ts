import { redirect } from '@sveltejs/kit';

export type UserBasicInfo = {
  email: string
  fullname: string
  username: string
}

export async function load({ cookies, fetch }) {
  let accessToken = cookies.get("access_token");
  if (!accessToken) {
    throw redirect(301, "/login");
  }

  const userList = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${accessToken}`,
    }
  });

  if (!userList.ok) {
    let a = await userList.json();
    if (userList.status < 500) {
      return {
        success: false,
        message: a.detail
      }
    }
  }

  let reponseData: UserBasicInfo[] = await userList.json()
  return { user_list: reponseData };
}
