package db

import (
	"fmt"
	"github.com/vladimish/ranhb/internal/db/models"
)

func (d *DataBase) GetTeachers(lastname string) ([]string, error) {
	var data []string
	query := fmt.Sprintf("SELECT teacher FROM ranh.teachers WHERE teacher LIKE '%%%s%%';", lastname)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}

		data = append(data, s)
	}

	return data, nil
}

func (d *DataBase) IsTeacherExists(teacher string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(teacher) FROM ranh.teachers WHERE teacher='%s';", teacher)
	var i int
	err := d.db.QueryRow(query).Scan(&i)
	if err != nil {
		return false, err
	}
	if i >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (d *DataBase) GetTeachersLessons(name, group string, day, month int) ([]models.TT, error) {
	query := fmt.Sprintf("SELECT tt.time, tt.subject, tt.classroom, tt.groups, tt.subgroup, tt.day_of_week FROM ranh.tt INNER JOIN ranh.teachers ON ranh.tt.teacher = ranh.teachers.id WHERE ranh.teachers.teacher='%s' AND ranh.tt.day='%d' AND ranh.tt.month='%d' AND ranh.tt.groups LIKE '%s';", name, day, month, group)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	var tts []models.TT
	for rows.Next() {
		var tt models.TT
		rows.Scan(&tt.Time, &tt.Subject, &tt.Classroom, &tt.Groups, &tt.Subgroup, &tt.Day_of_week)
		tt.Day = fmt.Sprintf("%d", day)
		tt.Month = fmt.Sprintf("%d", month)
		if group != "%" {
			tt.Groups = group
		}
		tt.Teacher = name
		tts = append(tts, tt)
	}

	return tts, nil
}
