INSERT INTO bot_stop (
  stop_time,
  start_time_id,
  wallet_balance
) SELECT
  NOW(),
  start.id,
  $1
FROM
  bot_start start
WHERE
  start.api_key_id = $2
  AND NOT EXISTS (
    select * from bot_stop stop where stop.start_time_id = start.id
  );
