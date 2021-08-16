package db

import (
	"fmt"
	"log"
	"time"
)

func (d *DataBase) CheckPremiumStatus(userId int64) (bool, error) {
	rows, err := d.db.Query(fmt.Sprintf("SELECT end_date FROM ranh.premium WHERE `user_id`=%d;", userId))
	if err != nil {
		return false, err
	}

	for rows.Next() {
		var timestamp int64
		rows.Scan(&timestamp)
		date := time.Unix(timestamp, 0)
		if date.After(time.Now()) {
			return true, nil
		}
	}

	return false, nil
}

func (d *DataBase) AddPremium(userId int64, duration time.Duration) error {
	isPremium, err := d.CheckPremiumStatus(userId)
	if err != nil {
		return err
	}
	if !isPremium {
		query := fmt.Sprintf("INSERT INTO ranh.premium (user_id, begin_date, end_date) VALUES (%d, %d, %d);", userId, time.Now().Unix(), time.Now().Add(duration).Unix())
		result, err := d.db.Exec(query)
		if err != nil{
			return err
		}
		log.Println(result)
	}

	return nil
}
