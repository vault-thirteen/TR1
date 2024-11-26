package cm

type Roles struct {
	CanLogIn        bool
	CanRead         bool
	CanWriteMessage bool
	CanCreateThread bool
	IsModerator     bool `gorm:"-"`
	IsAdministrator bool `gorm:"-"`
}
