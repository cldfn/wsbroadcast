package app

type EnvConfig struct {
	ApiPort          int
	CorsAllowOrigins string
}

func NewEnvConfig(envs EnvProvider) (*EnvConfig, error) {
	conf := &EnvConfig{}

	conf.ApiPort = envs.Int("API_HTTP_PORT", 5999)
	conf.CorsAllowOrigins = envs.String("CORS_ALLOW_ORIGINS", "*")

	return conf, nil

}
