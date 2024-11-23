package cm

type Password struct {
	MetaData
	Id     int `gorm:"primarykey"`
	UserId int `gorm:"index"`
	Bytes  []byte
}
