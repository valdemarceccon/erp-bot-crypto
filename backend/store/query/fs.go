package query

import _ "embed"

var (
	//go:embed user/new.sql
	NewUser string
	//go:embed user/list.sql
	ListUser string
	//go:embed user/stop_bot.sql
	StopBotUser string
	//go:embed user/start_bot.sql
	StarBotUser string
	//go:embed user/save_closed_pnl.sql
	SaveClosedPnLUser string
	//go:embed user/get_username.sql
	GetUsernameUser string
	//go:embed user/get_id.sql
	GetUser string
	//go:embed user/exists_username_email.sql
	UserExistsUsernameEmail string
)

var (
	//go:embed api_key/new.sql
	NewApiKey string
	//go:embed api_key/list.sql
	ListApiKeys string
	//go:embed api_key/from_user.sql
	GetApiFromUser string
	//go:embed api_key/get.sql
	GetApiKey string
	//go:embed api_key/save.sql
	SaveApiKey string
	//go:embed api_key/active.sql
	ListActiveApiKey string
	//go:embed api_key/bot_start_stop.sql
	BotStartStopApiKey string
)

var (
	//go:embed role/from_user.sql
	PermissionFromUser string
)
