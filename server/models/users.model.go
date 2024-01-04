package models

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const MIN_PASSWORD_LEN = 6
const MaxUsernameLen = 255

type User struct {
	gorm.Model
	Username             string `gorm:"column:username;unique;index" json:"username"`
	Password             string `gorm:"column:password" json:"password"`
	MinValidTokenVersion uint   `gorm:"column:min_valid_token_version;index"`
}

func (User) TableName() string {
	return "users"
}

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return errors.New("ParamError:missing username")
	}
	if len(username) > MaxUsernameLen {
		return fmt.Errorf("ParamError:username must be less than %d", MaxUsernameLen)
	}
	if match, _ := regexp.MatchString("^[a-zA-Z0-9-_.]+$", username); !match {
		return errors.New("ParamError:allowed characters in username are: a-z, A-Z, 0-9, ., _, -")
	}
	return nil
}

func (u *User) Create() error {
	if len(u.Password) < MIN_PASSWORD_LEN {
		return errors.New("ParamError:password too short")
	}
	if err := ValidateUsername(u.Username); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("PasswordError:%w", err)
	}
	u.Password = string(hash)
	if err := db.Create(u).Error; err != nil {
		return fmt.Errorf("DBError:%w", err)
	}
	return nil
}

func (u *User) GetByUsername() error {
	if err := ValidateUsername(u.Username); err != nil {
		return err
	}
	r := db.Where("username = ?", u.Username).First(u)
	if r.Error != nil {
		return fmt.Errorf("DBError:%w", r.Error)
	}
	if u.ID == 0 {
		return errors.New("ParamError:user not found")
	}
	return nil
}

func (u *User) GetById() error {
	r := db.First(u, u.ID)
	if r.Error != nil {
		return fmt.Errorf("DBError:%w", r.Error)
	}
	if u.ID == 0 {
		return errors.New("ParamError:user not found")
	}
	return nil
}
