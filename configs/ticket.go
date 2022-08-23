package configs

var TicketApp *TicketScheme

type TicketScheme struct {
	Web   WebConfig   `mapstructure:",squash"`
	Mysql MysqlConfig `mapstructure:",squash"`
}
