package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/michee/e-commerce/database"
	"github.com/michee/e-commerce/utils"
	"gorm.io/gorm"
)

var Db *gorm.DB

type User struct {
	UserId            string    `gorm:"primary_key;column:user_id"`
	Name              string    `gorm:"column:name" json:"name"`
	Username          string    `gorm:"column:username" json:"username"`
	Email             string    `gorm:"unique;column:email" json:"email"`
	Password          string    `gorm:"column:password" json:"password"`
	IsAdmin           bool      `gorm:"column:is_admin" json:"is_admin"`
	EmailVerified     bool      `gorm:"column:email_verified" json:"email_verified"`
	TokenJwt          string    `gorm:"column:token_jwt" json:"token_jwt"`
	VerificationToken string    `gorm:"column:verification_token" json:"verification_token"`
	TokenPassword     string    `gorm:"column:token_password" json:"token_password"`
	ResetTokenExpiry  time.Time `gorm:"column:reset_token_expiry" json:"reset_token_expiry"`
}

func init() {
	database.ConnectDB()
	Db = database.GetDB()
	// Db.Migrator().DropTable(&User{})
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
	u.UserId = uuid.New().String()
	admin := u.IsAdmin
	u.IsAdmin = admin
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
			return nil, nil
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
	u.TokenJwt = ""
	return Db.Save(&u).Error
}

func FindUserByToken(token string) (*User, error) {
	var u User
	if err := Db.Where("verification_token=?", token).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &u, nil
}

func FindUserPasswordToken(token string) (*User, error) {
	var u User
	if err := Db.Where("token_password = ?", token).First(&u).Error; err != nil {
		fmt.Println("Error finding user by token:", err)
		return nil, fmt.Errorf("token password not found")
	}

	// Vérifiez si le token a expiré
	if time.Now().After(u.ResetTokenExpiry) {
		return nil, fmt.Errorf("reset token has expired")
	}
	return &u, nil
}

func (u *User) Verify() error {
	u.EmailVerified = true
	u.VerificationToken = ""
	if err := Db.Save(&u).Error; err != nil {
		return fmt.Errorf("unable to verify user")
	}

	return nil
}

func (u *User) GeneratePasswordToken() error {
	token := uuid.New().String()
	u.TokenPassword = token
	u.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
	fmt.Println("Generated token:", token)
	fmt.Println("Token expiry time:", u.ResetTokenExpiry)
	return Db.Save(u).Error
}

func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	err = Db.Transaction(func(tx *gorm.DB) error {
		u.Password = hashedPassword
		u.TokenPassword = ""
		u.ResetTokenExpiry = time.Time{}

		if err := tx.Save(&u).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("New hased Password: ", hashedPassword)
	return nil
}
