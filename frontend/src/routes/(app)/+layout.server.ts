import { userTokenStore } from '$lib/stores.js'
import { error, redirect } from '@sveltejs/kit';
const BACKEND_HOST = process.env.BACKEND_HOST ? process.env.BACKEND_HOST : ""
export async function load({ params, cookies }) {

    const access_token = cookies.get("access_token",);
    if (!access_token) {
        throw redirect(301, "/login")
    }

    const userInfoEndpoint = new URL("/api/users/me", BACKEND_HOST)

    const user_data = await fetch(userInfoEndpoint, {
        headers: {
            "Authorization": `Bearer ${access_token}`,
        }
    });

    if (!user_data.ok) {
        cookies.delete("access_token");
        let error = await user_data.json();
        throw error(user_data.status, error.message);
    }

    let user_json = await user_data.json();

    return {
        access_token: access_token,
        email: user_json.email,
        name: user_json.name,
        username: user_json.username
    }
}
