package db

import "fmt"

func (d *DataBase) SaveCallback(chatId int64, messageId, day, month int) error {
	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO ranh.callback_messages (`message_id`, `chat_id`, `day`, `month`) VALUES (%d, %d, %d, %d);", messageId, chatId, day, month))
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBase) SaveTeacherCallback(chatId int64, messageId int, day int, month int, teacher string) error {
	_, err := d.db.Exec(fmt.Sprintf("INSERT INTO ranh.callback_messages (`message_id`, `chat_id`, `day`, `month`, `teacher`) VALUES (%d, %d, %d, %d, '%s');", messageId, chatId, day, month, teacher))
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBase) UpdateCallback(chatId int64, messageId, day, month int) error {
	_, err := d.db.Exec(fmt.Sprintf("UPDATE ranh.callback_messages SET `day`='%d', `month`='%d' WHERE `message_id`='%d' AND `chat_id`='%d';", day, month, messageId, chatId))
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBase) GetCallback(chatId int64, messageId int) (day int, month int, err error) {
	err = d.db.QueryRow(fmt.Sprintf("SELECT day, month FROM ranh.callback_messages WHERE message_id=%d AND chat_id=%d;", messageId, chatId)).Scan(&day, &month)
	if err != nil {
		return 0, 0, err
	}

	return day, month, nil
}


func (d *DataBase) GetTeacherCallback(chatId int64, messageId int) (day int, month int, teacher string, err error) {
	err = d.db.QueryRow(fmt.Sprintf("SELECT day, month, teacher FROM ranh.callback_messages WHERE message_id=%d AND chat_id=%d;", messageId, chatId)).Scan(&day, &month, &teacher)
	if err != nil {
		return 0, 0, "" , err
	}

	return day, month, teacher,nil
}