package project

import "github.com/gin-gonic/gin"

// RegisterRouter
func RegisterRouter(r *gin.RouterGroup) {

	// Routes
	r.GET("/projects", FindProjects)
	r.GET("/projects/:id", FindProject)
	r.POST("/projects", CreateProject)
	r.PATCH("/projects/:id", UpdateProject)
	r.DELETE("/projects/:id", DeleteProject)

}
