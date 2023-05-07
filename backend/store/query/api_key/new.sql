INSERT INTO api_key (
    user_id
    ,api_key_name
    ,exchange
    ,api_key
    ,api_secret
    ,status
    ,created_at
    ,updated_at
    ,deleted_at
  )
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), NULL)
RETURNING id;
