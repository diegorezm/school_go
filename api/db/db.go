package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBInfo struct {
	Host     string
	Port     int
	User     string
}

var Info = DBInfo{
	Host:     "localhost",
	Port:     3306,
	User:     "usuario",
}

type Databse struct {
	Connection *sql.DB
}

func NewDatabase(password, dbname string) *Databse {
	sqlinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", Info.User, password, Info.Host, Info.Port, dbname)
	db, err := sql.Open("mysql", sqlinfo)
	if err != nil {
		log.Fatalf("Could not load the databse. Erro: %s", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)
	return &Databse{Connection: db}
}
