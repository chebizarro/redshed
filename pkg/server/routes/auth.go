package routes

import (
	"github.com/chebizarro/redshed/internal/handlers/auth"
	"github.com/chebizarro/redshed/internal/orm"
	"github.com/chebizarro/redshed/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Auth routes
func Auth(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
	provider := string(utils.ProjectContextKeys.ProviderCtxKey)
	// OAuth handlers
	g := r.Group(cfg.VersionedEndpoint("/auth"))
	g.POST("/:"+provider+"/exchange", auth.Callback(cfg, orm))
	// g.GET(:"+provider+"/refresh", auth.Refresh(cfg, orm))
	return nil
}
