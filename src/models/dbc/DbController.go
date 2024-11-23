package dbc

import "gorm.io/gorm"

type DbController struct {
	db *gorm.DB
}

func NewDbController(db *gorm.DB) *DbController {
	return &DbController{db: db}
}
