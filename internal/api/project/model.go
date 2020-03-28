package project

import "github.com/jinzhu/gorm"

type Project struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
}

func SetupModel(db *gorm.DB) {
	db.AutoMigrate(&Project{})
}