SELECT
  s1.id
  ,api_key_id
  ,start_time
  ,s1.wallet_balance
  ,s2.stop_time
  ,s2.wallet_balance
FROM
  bot_start s1
INNER JOIN api_key ak
  ON s1.api_key_id = ak.id
LEFT JOIN
    bot_stop s2
  ON s1.id = s2.start_time_id
WHERE
  user_id = $1;
