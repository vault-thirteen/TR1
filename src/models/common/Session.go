package cm

import (
	"net"
)

type Session struct {
	MetaData
	Id       int `gorm:"primarykey"`
	UserId   int `gorm:"uniqueIndex"`
	UserIPAB net.IP
}
