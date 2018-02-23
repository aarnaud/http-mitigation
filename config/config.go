package config

var (
	Config *ServiceConfig
)

//ServerConfig is the server config struct
type ServiceConfig struct {
	HTTPPort      int
	CookieName    string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	Threshold1       int64
	Threshold2       int64
}
