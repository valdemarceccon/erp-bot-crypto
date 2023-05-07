SELECT id
  ,email
  ,username
  ,fullname
  ,hashed_password
FROM users
WHERE deleted_at IS NULL
ORDER BY id;
