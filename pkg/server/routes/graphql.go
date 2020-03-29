package routes

import (
	"github.com/chebizarro/redshed/internal/handlers"
	auth "github.com/chebizarro/redshed/internal/handlers/auth/middleware"
	"github.com/chebizarro/redshed/internal/logger"
	"github.com/chebizarro/redshed/internal/orm"
	"github.com/chebizarro/redshed/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GraphQL routes
func GraphQL(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
	// GraphQL paths
	gqlPath := cfg.VersionedEndpoint(cfg.GraphQL.Path)
	pgqlPath := cfg.GraphQL.PlaygroundPath
	g := r.Group(gqlPath)

	// GraphQL handler
	g.POST("", auth.Middleware(g.BasePath(), cfg, orm), handlers.GraphqlHandler(orm, &cfg.GraphQL))
	logger.Info("GraphQL @ ", gqlPath)
	// Playground handler
	if cfg.GraphQL.IsPlaygroundEnabled {
		logger.Info("GraphQL Playground @ ", g.BasePath()+pgqlPath)
		g.GET(pgqlPath, handlers.PlaygroundHandler(g.BasePath()))
	}

	return nil
}
