package configs

var TicketApp *TicketScheme

type TicketScheme struct {
	Web         WebConfig        `mapstructure:",squash"`
	Credentials Credentials      `mapstructure:",squash"`
	PostgreSQL  PostgreSQLConfig `mapstructure:",squash"`
	Language    LanguageConfig   `mapstructure:",squash"`
}
