package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

const (
	SCOPE_USER_PROTOTYPE_SELF       = "user.prototype.self"
	SCOPE_USER_PROTOTYPE_MANAGEMENT = "user.prototype.management"
)

func New(cfg *config.Config) *gin.Engine {
	// TODO scope register to scopes-service

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.Env.TrustProxy})

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg, SCOPE_USER_PROTOTYPE_MANAGEMENT)
	// users
	features.LoadUsers(engine, cfg, SCOPE_USER_PROTOTYPE_SELF)

	return engine
}
