package cm

import "time"

type MetaData struct {
	// Fields for GORM.
	CreatedAt time.Time `json:"toc"`
	UpdatedAt time.Time `json:"tou"`
}
