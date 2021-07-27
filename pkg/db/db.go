package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

// GetAllDistinctField Get all values of a specific field.
//
// Excepted arguments:
//
// 0 - name of the resulting field,
//
// 1 - name of the table,
func (d *DataBase) GetAllDistinctField(args ...string) ([]string, error) {
	if len(args) != 2 {
		panic("Expected 1 argument\n " +
			"0 - name of the resulting field,\n" +
			"1 - name of the table,\n" +
			"received: " + string(len(args)))
	}

	query := "SELECT DISTINCT `" + args[0] + "` FROM ranh." + args[1] + ";"
	return d.getAllFieldTemplate(query)
}

func (d *DataBase) getAllFieldTemplate(query string) ([]string, error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var tempGroup string
		err := rows.Scan(&tempGroup)
		if err != nil {
			return result, err
		}

		result = append(result, tempGroup)
	}

	return result, nil
}

func (d *DataBase) AddUser(id int, group string) error {
	query := fmt.Sprintf("INSERT INTO ranh.users (userId, primaryGroup) VALUES (%d, %s)", id, group)
	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
