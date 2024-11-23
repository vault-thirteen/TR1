package cm

import "time"

type MetaData struct {
	// Fields for GORM.
	CreatedAt time.Time
	UpdatedAt time.Time
}
