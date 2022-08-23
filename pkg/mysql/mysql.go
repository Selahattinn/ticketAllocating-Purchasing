package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	initTicketTable = `
	create table if not exists ticket (
    id bigint auto_increment primary key,
    name text not null,
    description text not null,
    quantity int  not null,
    constraint ticket_id_uindex
        unique (id)
);	
`
	initPurchaseTable = `
	create table if not exists purchase (
    id bigint auto_increment primary key,
    user_id text not null,
    quantity int  not null,
    constraint user_id_uindex
        unique (id)
);	
`
)

type Config struct {
	URL string
}

type IMysqlInstance interface {
	Database() *sql.DB
}

type mysqlInstance struct {
	database *sql.DB
}

func InitMysql(config Config) (IMysqlInstance, error) {
	db, err := sql.Open("mysql", config.URL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("creating tables error")
	}

	return &mysqlInstance{
		database: db,
	}, nil
}

func (p *mysqlInstance) Database() *sql.DB {
	return p.database
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(initTicketTable)
	if err != nil {
		return err
	}

	_, err = db.Exec(initPurchaseTable)
	if err != nil {
		return err
	}

	return nil
}
