package query

import _ "embed"

var (
	// go:embed api_key_new.sql
	NewApiKey string
	// go:embed api_key_list.sql
	ListApiKeys string
	// go:embed get_users.sql
	ListUsers string
	// go:embed api_key_from_user.sql
	GetApiFromUser string
)
