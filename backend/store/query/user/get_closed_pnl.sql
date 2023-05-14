select
    symbol,
    orderId,
    execType,
    closedPnl,
    createdTime,
    updatedTime
from
  closed_pnl
where
  user_id = $1
  and api_key_id = $2
  and (createdTime >= $3 and createdTime <= $4)
order by createdTime;
