INSERT INTO closed_pnl (
  user_id,
  api_key_id,
  symbol,
  orderId,
  execType,
  closedPnl,
  createdTime,
  updatedTime,
  created_at,
  updated_at
) VALUES (
  $1,    -- user_id
  $2,    -- api_key_id
  $3,    -- symbol
  $4,    -- orderId
  $5,    -- execType
  $6,    -- closedPnl
  $7,    -- createdTime
  $8,    -- updatedTime
  now(), -- created_at
  now()  -- updated_at
);
