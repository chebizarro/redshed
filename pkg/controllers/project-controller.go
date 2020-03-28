package controllers

import (
	"database/sql"
	"fmt"
	"github.com/chebizarro/redshed/pkg/models"
	"github.com/chebizarro/redshed/pkg/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"errors"
)

type ProjectController struct {
	projectRepository *repositories.ProjectRepository
}

func (c* ProjectController) Init(db * sql.DB) {
	c.projectRepository = &repositories.ProjectRepository{}
	c.projectRepository.Init(db)
}

func (c* ProjectController) CreateProject(ctx* gin.Context) {
	var project models.Project
	ctx.BindJSON(&project)
	
	if len(project.Title) == 0 {
		ctx.JSON(400, gin.H{"error":"title should not be empty"})
		return
	}
	
	createdProject, err := c.projectRepository.CreateProject(project)
	if err != nil {
		log.Printf("Error: %v\n", err)
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"project": createdProject,
	})

}

// GetProjects method
func (c *ProjectController) GetProjects(ctx *gin.Context) {
	projects := []models.Project{}
	var err error
	query := ctx.Request.URL.Query()
	pendings := query["pending"]
	froms := query["from"]
	tos := query["to"]

	useridi, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "userid not found in request context",
		})
		return
	}
	userid := useridi.(string)

	if len(pendings) > 0 {
		log.Printf("invoking GetPendingProjects for user %s\n", userid)
		projects, err = c.projectRepository.GetPendingProjects(userid)
	} else if len(froms) > 0 && len(tos) > 0 {
		from := froms[0]
		to := tos[0]
		log.Printf("invoking GetProjectsByDateRange for user=%s with from=%s, to=%s\n", userid, from, to)
		projects, err = c.projectRepository.GetProjectsByDateRange(userid, from, to)
	} else {
		errorMessage := "Invalid query parameters. Either pending or from/to should be provided"
		log.Println(errorMessage)
		err = errors.New(errorMessage)
	}

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"projects": projects,
	})
}

// GetProjectByID method
func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	useridi, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "userid not found in request context",
		})
		return
	}
	userid := useridi.(string)

	task, err := c.projectRepository.GetProjectByID(userid, id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"task": task,
	})
}

// UpdateProjectForID method
func (c *ProjectController) UpdateProjectForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	useridi, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "userid not found in request context",
		})
		return
	}
	userid := useridi.(string)

	existingProject, err := c.projectRepository.GetProjectByID(userid, id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.BindJSON(&existingProject)
	err = c.projectRepository.UpdateProject(userid, id, existingProject)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set the user tags
	//c.userTagsRepository.SetUserTags(userid, existingProject.Tags)

	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("%d updated", id),
	})
}

// DeleteProjectForID method
func (c *ProjectController) DeleteProjectForID(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": fmt.Sprintf("%s is not a valid number", idstr),
		})
		return
	}

	useridi, exists := ctx.Get("userid")
	if !exists {
		ctx.JSON(400, gin.H{
			"error": "userid not found in request context",
		})
		return
	}
	userid := useridi.(string)

	err = c.projectRepository.DeleteProject(userid, id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("%d deleted", id),
	})
}
