package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CreateProjectInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type UpdateProjectInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// GET /projects
// Find all projects
func FindProjects(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var projects []Project
	db.Find(&projects)

	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// GET /projects/:id
// Find a project
func FindProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var project Project
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// POST /projects
// Create new project
func CreateProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create project
	project := Project{Name: input.Author}
	db.Create(&project)

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// PATCH /projects/:id
// Update a project
func UpdateProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var project Project
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&project).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": project})
}

// DELETE /projects/:id
// Delete a project
func DeleteProject(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var project Project
	if err := db.Where("id = ?", c.Param("id")).First(&project).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&project)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

