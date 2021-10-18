package db

import (
	"fmt"
	"github.com/telf01/ranhb/pkg/users/models"
	"time"
)

func (d *DataBase) AddUser(u *models.User) error {
	query := fmt.Sprintf("INSERT INTO ranh.users (userId, primary_group, last_action, last_action_value, is_privacy_accepted, is_terms_of_use_accepted, register_time, last_action_time, actions_amount) VALUES (%d, \"%s\", \"%s\", %d, %t, %t, %d, %d, %d) ON DUPLICATE KEY UPDATE userId=%d, primary_group=\"%s\", last_action=\"%s\", last_action_value=%d, is_privacy_accepted=%t, is_terms_of_use_accepted=%t, last_action_time=%d, actions_amount=%d;", u.Id, u.Group, u.LastAction, u.LastActionValue, u.IsPrivacyAccepted, u.IsTermsOfUseAccepted, time.Now().Unix(), time.Now().Unix(), 0, u.Id, u.Group, u.LastAction, u.LastActionValue, u.IsPrivacyAccepted, u.IsTermsOfUseAccepted, time.Now().Unix(), u.ActionsAmount)
	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) GetUser(id int64, u **models.User) error {
	query := "SELECT * FROM ranh.users WHERE userId =" + fmt.Sprintf("%d", id) + ";"
	row, err := d.db.Query(query)
	if err != nil {
		return err
	}

	if *u == nil {
		*u = &models.User{}
	}

	for row.Next() {
		err := row.Scan(&(*u).Id, &(*u).Group, &(*u).LastAction, &(*u).LastActionValue, &(*u).IsPrivacyAccepted, &(*u).IsTermsOfUseAccepted, &(*u).RegisterTime, &(*u).LastActionTime, &(*u).ActionsAmount)
		if err != nil {
			return err
		}
	}

	if (*u).Id == 0 {
		(*u).Id = id
		(*u).LastAction = "start"
	}

	err = d.UpdateUserTimestamp(id)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataBase) UpdateUserTimestamp(id int64) error {
	query := fmt.Sprintf("UPDATE ranh.users SET last_action_time=%d, actions_amount=actions_amount+1 WHERE userId=%d;", time.Now().Unix(), id)
	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
