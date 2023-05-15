import { redirect } from '@sveltejs/kit';

type Commission = {
  date: Date
  profit: number
  fee: number
  high_mark: number
  balance: number
}

type CommissionResp = {
  start: Date
  stop: Date | null
  username: string
  commissions: Commission[]
};

export async function load({ cookies, fetch }) {
  const token = cookies.get("access_token");

  if (!token) {
    redirect(301, "/");
  }

  let api_keys_resp = await fetch(`http://${process.env.BACKEND_PRIVATE_HOST}/user/comission/`, {
    method: "GET",
    headers: {
      "Authorization": `Bearer ${token}`,
    }
  });

  if (!api_keys_resp.ok) {
    let a = await api_keys_resp.json();
    return {
      success: false,
      message: a.message
    }
  }

  let commissions: CommissionResp[] = [];

  for (let r of await api_keys_resp.json()) {
    let botrun: CommissionResp = {
      start: new Date(r.start),
      stop: r.stop ? new Date(r.stop) : null,
      username: r.username,
      commissions: []
    }

    for (let c of r.commissions) {
      botrun.commissions.push({
        balance: c.balance,
        profit: c.profit,
        fee: c.fee,
        high_mark: c.high_mark,
        date: new Date(c.date)
      })
    }

    commissions.push(botrun)
  }

  return {
    success: true,
    commissions: commissions
  }

}
