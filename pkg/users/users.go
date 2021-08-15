package users

import (
	_ "database/sql"
	"github.com/telf01/ranhb/pkg/db"
	"github.com/telf01/ranhb/pkg/users/models"
)

type User struct {
	U  *models.User
	db *db.DataBase
}

func NewUser(id int64, base *db.DataBase) *User {
	return &User{U: &models.User{Id: id}, db: base}
}

func Get(id int64, base *db.DataBase) (*User, error) {
	var user User
	err := base.GetUser(id, &user.U)
	if err != nil {
		return nil, err
	}
	user.db = base
	return &user, nil
}

func (u *User) Save() error {
	err := u.db.AddUser(u.U)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Init(id int64, db *db.DataBase) error {
	u = NewUser(id, db)
	u.U.LastAction = "start"
	err := u.Save()
	if err != nil{
		return err
	}

	return nil
}
