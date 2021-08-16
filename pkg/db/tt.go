package db

import (
	"fmt"
	"github.com/telf01/ranhb/pkg/db/models"
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

	query := fmt.Sprintf("SELECT DISTINCT `%s` FROM ranh.%s ORDER BY `%s` LIMIT %s, %s;", args[0], args[1], args[0], args[2], args[3])
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

func (d *DataBase) GetSpecificTt(group string, day int, month int) ([]models.TT, error) {
	query := fmt.Sprintf("SELECT * FROM ranh.tt WHERE `groups`='%s' AND `day`='%d' AND `month`='%d' ORDER BY `time`;", group, day, month)
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

	var r models.TT
	var id int

	var result []models.TT

	for rows.Next() {
		err := rows.Scan(&id, &r.Day, &r.Month, &r.Day_of_week, &r.Time, &r.Amount, &r.Groups, &r.Subject_type, &r.Subject, &r.Rank, &r.Classroom, &r.Teacher, &r.Fuck_key, &r.Subgroup)
		if err != nil {
			return nil, err
		}

		err = d.db.QueryRow(fmt.Sprintf("SELECT teacher FROM ranh.teachers WHERE `id`=%s;", r.Teacher)).Scan(&r.Teacher)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return result, nil
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
		var field string
		err := rows.Scan(&field)
		if err != nil {
			return result, err
		}

		result = append(result, field)
	}

	return result, nil
}
