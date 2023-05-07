SELECT id
  ,user_id
  ,api_key_name
  ,exchange
  ,api_key
  ,api_secret
  ,status
FROM api_key
WHERE id = $1
  AND user_id = $2
ORDER BY
  id,
  user_id;
