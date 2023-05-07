INSERT INTO users(
  email,
  username,
  fullname,
  hashed_password,
  created_at,
  updated_at
)
VALUES ($1,$2,$3,$4,now(),now())
RETURNING id;
