package user

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Username	string	`json:"username"`
	Email		*string	`json:"email" gorm:"unique;not null"`
	Password	string	`json:"password" form:"password" validate:"required" gorm:"password"`
	Token 		string
	Role		string	`json:"role" form:"role"`
	Active		bool
}