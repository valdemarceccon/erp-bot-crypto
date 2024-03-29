import { error, redirect } from '@sveltejs/kit';

export async function load({ cookies, fetch }) {

    const access_token = cookies.get("access_token",);
    if (!access_token) {
        throw redirect(301, "/login")
    }
    try {
        const user_data = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/me`, {
            headers: {
                "Authorization": `Bearer ${access_token}`,
            }
        });
        if (!user_data.ok) {
            cookies.delete("access_token");
            let error_resp = await user_data.json();
            throw error(user_data.status, error_resp.message);
        }

        let user_json: {
            success: boolean,
            access_token: string,
            email: string,
            fullname: string,
            username: string,
            permissions: {
                name: string
            }[]
        } = { ...await user_data.json(), access_token: access_token, success: true };

        return {
            success: true,
            access_token: access_token,
            email: user_json.email,
            name: user_json.fullname,
            username: user_json.username,
            permissions: user_json.permissions
        }
    } catch (e: any) {
        return {
            success: false,
            message: e.message
        }
    }


}
