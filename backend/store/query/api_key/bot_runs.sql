select
  *
from
  bot_start s1,
  bot_stop s2,
  api_key a
where
    s1.api_key_id = a.id
and s2.start_time_id = s1.id
and api_key_id = $1;
