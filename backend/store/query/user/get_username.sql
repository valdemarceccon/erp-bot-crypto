select
  id,
  email,
  username,
  fullname,
  hashed_password
from users
where username = $1
  and deleted_at is null;
