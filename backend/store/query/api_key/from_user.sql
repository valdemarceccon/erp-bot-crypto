SELECT id
  ,user_id
  ,api_key_name
  ,exchange
  ,api_key
  ,api_secret
  ,status
FROM api_key
WHERE deleted_at IS NULL
  AND user_id = $1
ORDER BY id,
  user_id;
