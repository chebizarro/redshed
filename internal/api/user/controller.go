package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type CreateUserInput struct {
	Username  string `json:"username" binding:"required"`
}

type UpdateUserInput struct {
	Username  string `json:"username"`
}

func SetupModel(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

// GET /users
// Find all users
func FindUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []User
	db.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GET /users/:id
// Find a user
func FindUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var user User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// @Title User Registration
// @Summary Allows a new user to register for the service
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200
// @Failure 500
// @Router /api/user/register [post]
func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user := User{Username: input.Username}
	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// PATCH /users/:id
// Update a user
func UpdateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var user User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// DELETE /users/:id
// Delete a user
func DeleteUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Get model if exist
	var user User
	if err := db.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

