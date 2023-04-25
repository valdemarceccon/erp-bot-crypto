import { error, redirect } from '@sveltejs/kit';
import { BACKEND_PRIVATE_HOST } from '$env/static/private';

export async function load({ cookies, fetch }) {

    const access_token = cookies.get("access_token",);
    if (!access_token) {
        throw redirect(301, "/login")
    }
    try {
        const user_data = await fetch(`http://${BACKEND_PRIVATE_HOST}/users/me`, {
            headers: {
                "Authorization": `Bearer ${access_token}`,
            }
        });
        if (!user_data.ok) {
            cookies.delete("access_token");
            let error_resp = await user_data.json();
            throw error(user_data.status, error_resp.message);
        }

        let user_json = await user_data.json();

        return {
            success: true,
            access_token: access_token,
            email: user_json.email,
            name: user_json.name,
            username: user_json.username
        }
    } catch (e: any) {
        return {
            success: false,
            message: e.message
        }
    }


}
