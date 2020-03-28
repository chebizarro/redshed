package zone

import (
	"github.com/jinzhu/gorm"
	"github.com/paulmach/orb"
)

type Zone struct {
	ID     uint   `json:"id" gorm:"primary_key"`
	Name   string `json:"name"`
	Geometry orb.Polygon `json:"geometry" gorm:"type:byte[]"`
}

func SetupModel(db *gorm.DB) {
	db.AutoMigrate(&Zone{})
}