package db

import (
	"fmt"
	"github.com/vladimish/yookassa-go-sdk"
	"log"
	"time"
)

func (d *DataBase) SavePayment(payment yookassa.Payment, userID int64) error {
	query := fmt.Sprintf("INSERT INTO ranh.payments (payment_id, payment_time, user) VALUES('%s', %d, %d);", payment.Id, time.Now().Unix(), userID)
	res, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}
