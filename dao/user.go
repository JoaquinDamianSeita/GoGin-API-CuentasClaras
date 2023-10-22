package dao

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `gorm:"column:id; primary_key; not null" json:"id"`
	Username string `gorm:"column:username; unique" json:"username"`
	Email    string `gorm:"column:email; unique" json:"email"`
	Password string `gorm:"column:password" json:"password"`
	BaseModel
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
