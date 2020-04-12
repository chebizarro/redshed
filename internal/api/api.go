package api

import (
	"github.com/casbin/casbin"
	_ "github.com/chebizarro/redshed/docs"
	"github.com/chebizarro/redshed/internal/api/project"
	"github.com/chebizarro/redshed/internal/api/task"
	"github.com/chebizarro/redshed/internal/api/user"
	"github.com/chebizarro/redshed/internal/api/zone"
	"github.com/chebizarro/redshed/internal/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func InitRouter() *gin.Engine {

	router := gin.Default()

	router.Use(CORS())

	setUpConfig(router)
	setUpRouter(router)

	return router
}

func setUpConfig(router *gin.Engine) {

	router.Static("/assets", "./web/dist")
	router.StaticFile("/app.html", "./web/dist/index.html")

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Connect to database
	db := SetupModels()

	// Provide db to controllers
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	e := casbin.NewEnforcer("config/authz/model.conf", "config/authz/policy.csv")
	router.Use(middleware.NewAuthorizer(e))
}

func setUpRouter(router *gin.Engine) {
	api := router.Group("/api")
	{
		user.RegisterRouter(api.Group("/auth"))
		project.RegisterRouter(api.Group("/project"))
	}
}

func SetupModels() *gorm.DB {
	db, err := gorm.Open("sqlite3", "web.db")

	if err != nil {
		panic("Failed to connect to database!" + err.Error())
	}
	defer db.Close()

	project.SetupModel(db)
	task.SetupModel(db)
	user.SetupModel(db)
	zone.SetupModel(db)

	return db
}
