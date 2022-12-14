package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/Selahattinn/ticketAllocating-Purchasing/configs"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/logging"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/viperconfig"
)

func boot(logger *logrus.Logger) (*application, error) {
	postgresInstance, postgresErr := initMysql()
	if postgresErr != nil {
		return nil, postgresErr
	}

	return &application{
		logger:          logger,
		postgreInstance: postgresInstance,
	}, nil
}

func initConfig() error {
	path := "."
	if envConfigPath := os.Getenv("CONFIG_FILE_PATH"); envConfigPath != "" {
		path = envConfigPath
	}

	file := ".env"
	if envConfigFile := os.Getenv("CONFIG_FILE_NAME"); envConfigFile != "" {
		file = envConfigFile
	}

	vc := viperconfig.Config{
		Path:     path,
		FileName: file,
		Env:      os.Getenv("ENV"),
	}

	c, err := viperconfig.Load(vc, configs.TicketScheme{})
	if err != nil {
		return err
	}

	config := c.(configs.TicketScheme)
	configs.TicketApp = &config

	return nil
}

func initLogger() *logrus.Logger {
	return logging.NewLogger(logging.Config{
		Service: logging.ServiceConfig{
			Env:     configs.TicketApp.Web.Env,
			AppName: configs.TicketApp.Web.AppName,
		},
	})
}

func initMysql() (mysql.IMysqlInstance, error) {
	return mysql.InitMysql(mysql.Config{
		URL: configs.TicketApp.Mysql.URL,
	})
}
