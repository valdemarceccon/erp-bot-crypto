SELECT
  1
FROM
  users
WHERE
    username = $1
  or email = $2
  and deleted_at is null;
