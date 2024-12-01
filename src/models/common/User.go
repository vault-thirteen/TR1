package cm

import (
	"net/mail"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
)

type User struct {
	MetaData
	Id                int        `json:"id" gorm:"primarykey"`
	Name              string     `json:"name" gorm:"uniqueIndex,size:255"`
	Email             string     `json:"email" gorm:"uniqueIndex,size:255"`
	Password          *Password  `json:"password,omitempty"`
	Session           *Session   `json:"session,omitempty"`
	Roles             Roles      `json:"roles" gorm:"embedded"`
	RegTime           time.Time  `json:"regTime"`
	LastBadActionTime *time.Time `json:"lastBadActionTime"`
	BanTime           *time.Time `json:"banTime"`
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
