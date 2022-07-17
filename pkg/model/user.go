package model

import (
	"github.com/grin-ch/grin-api/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	api.Model
	Username    string `json:"username" gorm:"username"`
	Password    string `json:"password" gorm:"password"`
	Email       string `json:"email" gorm:"email"`
	PhoneNumber string `json:"phoneNumber" gorm:"phoneNumber"`
}

func (p *Provider) CreateUser(user User) error {
	tx := p.db.Begin()
	tx.Create(&user)

	tx = tx.Where("username = ?", user.Username)
	if user.Email != "" {
		tx = tx.Or("email = ?", user.Email)
	}
	if user.PhoneNumber != "" {
		tx = tx.Or("phone_number = ?", user.PhoneNumber)
	}
	var count int64
	tx.Count(&count)
	if count > 1 {
		tx.Rollback()
		return status.Errorf(codes.AlreadyExists, "username, email or phoneNumber exits")
	}

	return tx.Commit().Error
}

func (p *Provider) PickUser(user User) (User, error) {
	u := User{}
	err := p.db.Model(u).Where(user).First(&u).Error
	return u, err
}
