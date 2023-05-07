SELECT
	p.permission_name
FROM
	user_roles ur
	inner join role_permission rp on
		ur.role_id = rp.role_id
	inner join permission p on
		rp.permission_id = p.id
WHERE
	ur.user_id = $1
ORDER BY
	id;
