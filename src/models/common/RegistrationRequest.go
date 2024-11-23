package cm

import "net"

type RegistrationRequest struct {
	MetaData
	Id int `gorm:"primarykey"`

	// Fields requested by a potential user.
	UserName     string `gorm:"uniqueIndex,size:255"`
	UserEmail    string `gorm:"uniqueIndex,size:255"`
	UserPassword string

	// System fields.
	RequestId          string `gorm:"uniqueIndex,size:255"`
	UserIPAB           net.IP
	CaptchaId          string
	VerificationCode   string
	IsReadyForApproval bool `gorm:"index"`
}
