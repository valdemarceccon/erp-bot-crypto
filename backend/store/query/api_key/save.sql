UPDATE api_key
SET api_key_name = $1
  ,exchange = $2
  ,api_key = $3
  ,api_secret = $4
  ,status = $5
WHERE
  id = $6;
