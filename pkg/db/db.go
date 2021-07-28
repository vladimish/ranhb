package db

import (
	"database/sql"
	"fmt"
)

type DataBase struct {
	db *sql.DB
}

func NewDataBase(db *sql.DB) *DataBase {
	return &DataBase{db: db}
}

func InitDBConncetion(login string, password string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ranh", login, password))
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}
