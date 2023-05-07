SELECT 	id,
        email,
        username,
        fullname,
        hashed_password
FROM users
WHERE id = $1
  AND deleted_at is null;
