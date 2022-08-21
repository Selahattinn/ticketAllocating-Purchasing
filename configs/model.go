package configs

type Credentials struct {
	AppID string `mapstructure:"APP_ID"`
}

type WebConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    string `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
}

type PostgreSQLConfig struct {
	URL      string `mapstructure:"POSTGRESQL_URL"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Database string `mapstructure:"POSTGRES_DB"`
}
