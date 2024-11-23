package cm

import (
	"net/mail"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
)

type User struct {
	MetaData
	Id                int    `gorm:"primarykey"`
	Name              string `gorm:"uniqueIndex,size:255"`
	Email             string `gorm:"uniqueIndex,size:255"`
	Password          *Password
	Roles             Roles `gorm:"embedded"`
	RegTime           time.Time
	LastBadActionTime *time.Time
	BanTime           *time.Time
}

func IsUserEmailValid(email string) (isValid bool) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func IsUserPasswordAllowed(password string) (isAllowed bool) {
	ok, _ := bpp.IsPasswordAllowed(password)
	return ok
}
