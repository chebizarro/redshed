package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter
func RegisterRouter(r *gin.RouterGroup) {

	r.POST("/register", CreateUser)

	r.POST("/login", Auth.LoginHandler)

	auth := r.Group("")
	auth.Use(Auth.MiddlewareFunc())
	{

		auth.GET("", FindUsers)

		auth.DELETE("/:id", DeleteUser)

		auth.PUT("/:id", UpdateUser)
	}
}
