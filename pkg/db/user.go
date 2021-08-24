package db

import (
	"fmt"
	"github.com/telf01/ranhb/pkg/users/models"
)

func (d *DataBase) AddUser(u *models.User) error {
	query := fmt.Sprintf("INSERT INTO ranh.users (userId, primary_group, last_action, last_action_value, is_privacy_accepted, is_terms_of_use_accepted) VALUES (%d, \"%s\", \"%s\", %d, %t, %t) ON DUPLICATE KEY UPDATE userId=%d, primary_group=\"%s\", last_action=\"%s\", last_action_value=%d, is_privacy_accepted=%t, is_terms_of_use_accepted=%t;", u.Id, u.Group, u.LastAction, u.LastActionValue, u.IsPrivacyAccepted, u.IsTermsOfUseAccepted, u.Id, u.Group, u.LastAction, u.LastActionValue, u.IsPrivacyAccepted, u.IsTermsOfUseAccepted)
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
		err := row.Scan(&(*u).Id, &(*u).Group, &(*u).LastAction, &(*u).LastActionValue, &(*u).IsPrivacyAccepted, &(*u).IsTermsOfUseAccepted)
		if err != nil {
			return err
		}
	}

	if (*u).Id == 0{
		(*u).Id = id
		(*u).LastAction = "start"
	}

	return nil
}
