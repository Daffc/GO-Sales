package model

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;default:auto_random()"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrUserNameRequired                   = errors.New("invalid user name")
	ErrUserEmailRequired                  = errors.New("invalid user email")
	ErrUserPasswordFormat                 = errors.New("password must have at least one uppercase letter, one lowercase letter, one numeric character, and one special character")
	ErrUserPasswordLenght                 = errors.New("password be at least 6 characters long")
	ErrUserPasswordFormatLowCase          = errors.New("the password must have at least one lowercase character")
	ErrUserPasswordFormatUpperCase        = errors.New("the password must have at least one uppercase character ")
	ErrUserPasswordFormatNumber           = errors.New("the password must have at least one numeric character")
	ErrUserPasswordFormatSpecialCharacter = errors.New("the password must have at least one special character")
)

func (u *User) ValidatePassword() error {
	if len(u.Password) < 6 {
		return ErrUserPasswordLenght
	}

	re := regexp.MustCompile(`(.*[a-z])+`)
	if !re.MatchString(u.Password) {
		return ErrUserPasswordFormatLowCase
	}

	re = regexp.MustCompile(`(.*[A-Z])+`)
	if !re.MatchString(u.Password) {
		return ErrUserPasswordFormatUpperCase
	}

	re = regexp.MustCompile(`(.*[1-9])+`)
	if !re.MatchString(u.Password) {
		return ErrUserPasswordFormatNumber
	}

	re = regexp.MustCompile(`(.*[-._!"\x60´'#%&,:;<>=@{}~\$\(\)\*\+\/\\\?\[\]\^\|])+`)
	if !re.MatchString(u.Password) {
		return ErrUserPasswordFormatSpecialCharacter
	}

	return nil
}

func (u *User) ValidateName() error {
	if len(u.Name) == 0 {
		return ErrUserNameRequired
	}

	return nil
}

func (u *User) ValidateEmail() error {
	re := regexp.MustCompile(`^([a-zA-Z0-9._-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z]+)+)$`)
	if !re.MatchString(u.Email) {
		return ErrUserEmailRequired
	}

	return nil
}

func (u *User) ValidateAll() error {

	if err := u.ValidateName(); err != nil {
		return err
	}

	if err := u.ValidateEmail(); err != nil {
		return err
	}
	if err := u.ValidatePassword(); err != nil {
		return err
	}

	return nil
}
