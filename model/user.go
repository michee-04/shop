package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/michee/e-commerce/database"
	"github.com/michee/e-commerce/utils"
	"gorm.io/gorm"
)

var Db *gorm.DB

type User struct {
	UserId            string    `gorm:"primary_key;column:user_id"`
	Username          string    `gorm:"column:username" json:"username"`
	Email             string    `gorm:"unique:column:email" json:"email"`
	EmailVerified     bool      `gorm:"column:email_verified" json:"email_verified"`
	Password          string    `gorm:"column:password" json:"password"`
	Token             string    `gorm:"column:token" json:"token"`
	VerificationToken string    `gorm:"column:verification_token" json:"verification_token"`
	ResetTokenExpiry  time.Time `gorm:"column:reset_token_expiry" json:"reset_token_expiry"`
}

// func (u *User) HexaDeximalId(tx *gorm.DB) (err error) {
// 	u.UserId = uuid.New().String()
// 	return
// }

func init() {
	database.ConnectDB()
	Db = database.GetDB()

	if Db != nil {
		err := Db.AutoMigrate(&User{})
		if err != nil {
			panic("Failed to migrate User model: " + err.Error())
		}
	} else {
		panic("DB connection is nil")
	}
}

func (u *User) CreateUser() *User {
	hashedPassword, _ := utils.HashPassword(u.Password)
	u.Password = hashedPassword
	Db.Create(u)
	return u
}

func GetUser() []User {
	var u []User
	Db.Find(&u)
	return u
}

func GetUserById(Id string) (*User, *gorm.DB) {
	var u User
	db := Db.Where("user_id=?", Id).First(&u)
	if db.Error != nil {
		return nil, db
	}

	return &u, db
}

func GetUserByEmail(email string) (*User, error) {
	var u User
	if err := Db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func DeleteUser(Id string) User {
	var u User
	Db.Where("user_id=?", Id).Delete(&u)
	return u
}

func (u *User) Logout() error {
	u.Token = ""
	return Db.Save(&u).Error
}

func FindUserByToken(token string) (*User, error) {
	var u User
	if err := Db.Where("verification_token=?", token).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &u, nil
}

func (u *User) Verify() error {
	u.EmailVerified = true
	u.VerificationToken = ""
	if err := Db.Save(&u).Error; err != nil {
		return fmt.Errorf("Unable to verify user")
	}

	return nil
}


func (u *User) CanLogin() bool{
	return u.EmailVerified
}

