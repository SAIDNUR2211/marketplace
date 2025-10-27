package configs

type Configs struct {
	AppParams      AppParams      `json:"app_params"`
	PostgresParams PostgresParams `json:"postgres_params"`
	AuthParams     AuthParams     `json:"auth_params"`
	RedisParams    RedisParams    `json:"redis_params"`
}
type AppParams struct {
	ServerURL  string `json:"server_url"`
	ServerName string `json:"server_name"`
	PortRun    string `json:"port_run"`
	GinMode    string `json:"gin_mode"`
}
type PostgresParams struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}
type AuthParams struct {
	AccessTokenTtlMinutes int `json:"access_token_ttl_minutes"`
	RefreshTokenTtlDays   int `json:"refresh_token_ttl_days"`
}
type RedisParams struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}
