SELECT api.id
  ,user_id
  ,u.username
  ,api_key_name
  ,exchange
  ,api_key
  ,api_secret
  ,status
FROM api_key api
  INNER JOIN users u ON u.id = api.user_id
WHERE api.deleted_at IS NULL
  AND u.deleted_at IS NULL
ORDER BY id,
  user_id;
