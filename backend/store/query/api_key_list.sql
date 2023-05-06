SELECT ak.id,
  user_id,
  u.username,
  api_key_name,
  exchange,
  api_key,
  api_secret,
  status
FROM api_key ak
  INNER JOIN users u ON u.id = ak.user_id
WHERE ak.deleted_at IS NULL
  AND u.deleted_at IS NULL
ORDER BY id,
  user_id;
