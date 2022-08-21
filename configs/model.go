package configs

import (
	"golang.org/x/text/language"
)

type Credentials struct {
	AppID  string `mapstructure:"APP_ID"`
	Secret string `mapstructure:"APP_SECRET"`
}

type WebConfig struct {
	AppName string `mapstructure:"APP_NAME"`
	Port    string `mapstructure:"PORT"`
	Env     string `mapstructure:"ENV"`
}

type PostgreSQLConfig struct {
	URL      string `mapstructure:"POSTGRESQL_URL"`
	Database string `mapstructure:"POSTGRESQL_DATABASE"`
}

type LanguageConfig struct {
	Default   language.Tag `mapstructure:"LANGUAGE_DEFAULT"`
	Languages []language.Tag
}
