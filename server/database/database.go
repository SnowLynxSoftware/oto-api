package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type AppDataSource struct {
	DB *sqlx.DB
}

func NewAppDataSource() *AppDataSource {
	return &AppDataSource{}
}

func (d *AppDataSource) Connect(dbConnectionString string) {
	db, err := sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}
	d.DB = db
	err = d.DB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database is connected")
}
