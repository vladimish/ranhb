package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// GetAllDistinctField Get all values of a specific field.
//
// Excepted arguments:
//
// 0 - name of the resulting field,
//
// 1 - name of the table,
//
// 2 - limit offset,
//
// 3 - page size.
func (d *DataBase) GetAllDistinctField(args ...string) ([]string, error) {
	if len(args) != 4 {
		panic("Expected 4 arguments\n " +
			"0 - name of the resulting field,\n" +
			"1 - name of the table,\n" +
			"2 - limit offset,\n" +
			"3 - page size,\n" +
			"received: " + fmt.Sprintf("%d", len(args)))
	}

	query := fmt.Sprintf("SELECT DISTINCT `%s` FROM ranh.%s LIMIT %s, %s;", args[0], args[1], args[2], args[3])
	return d.getAllFieldTemplate(query)
}

// GetAllDistinctFieldWhere Get all values of a specific field.
//
// Excepted arguments:
//
// 0 - name of the resulting field,
//
// 1 - name of the table,
//
// 2 - limit offset,
//
// 3 - page size,
//
// 4 - where name,
//
// 5 - where value.
func (d *DataBase) GetAllDistinctFieldWhere(args ...string) ([]string, error) {
	if len(args) != 6 {
		panic("Expected 6 arguments\n " +
			"0 - name of the resulting field,\n" +
			"1 - name of the table,\n" +
			"2 - limit offset,\n" +
			"3 - page size,\n" +
			"4 - where name,\n" +
			"5 - where value,\n" +
			"received: " + fmt.Sprintf("%d", len(args)))
	}

	query := fmt.Sprintf("SELECT DISTINCT `%s` FROM ranh.%s WHERE `%s`='%s' LIMIT %s, %s;", args[0], args[1], args[4], args[5], args[2], args[3])
	return d.getAllFieldTemplate(query)
}

func (d *DataBase) getAllFieldTemplate(query string) ([]string, error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}()

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
