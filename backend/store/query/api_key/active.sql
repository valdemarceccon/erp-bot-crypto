SELECT a.id
  ,a.user_id
  ,a.api_key_name
  ,a.exchange
  ,a.api_key
  ,a.api_secret
  ,a.status
FROM users u
inner join api_key a on
        u.id = a.user_id
    and a.deleted_at is null
    and u.deleted_at is null
WHERE
    (u.id = $1 OR $1 = 0)
and a.status in ($2, $3);
