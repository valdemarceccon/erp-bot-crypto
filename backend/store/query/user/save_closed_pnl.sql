INSERT INTO closed_pnl (
  user_id, api_key_id, symbol, orderId, side, qty, orderPrice,
  orderType, execType, closedSize, cumEntryValue, avgEntryPrice,
  cumExitValue, avgExitPrice, closedPnl, fillCount, leverage,
  createdTime, updatedTime, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,$13, $14, $15, $16, $17, $18, $19, now(), now()
);
