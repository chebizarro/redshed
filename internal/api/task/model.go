package task

import "github.com/jinzhu/gorm"

type Task struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
	Description string `json:"description"`
}

func SetupModel(db *gorm.DB) {
	db.AutoMigrate(&Task{})
}