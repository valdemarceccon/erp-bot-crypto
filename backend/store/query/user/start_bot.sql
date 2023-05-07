INSERT INTO bot_start(
  api_key_id,
  start_time,
  wallet_balance
)
values (
  $1, now(), $2
);
