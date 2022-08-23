package configs

const (
	productionEnv = "production"
)

type WebConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    string `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
}

type MysqlConfig struct {
	URL string `mapstructure:"MYSQL_URL"`
}

func (wc WebConfig) IsProductionEnv() bool {
	return wc.Env == productionEnv
}
