package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

const (
	SCOPE_USER_PROTOTYPE_SELF = "user.prototype.self"
	SCOPE_USER_PROTOTYPE      = "user.prototype"
)

func New(cfg config.ConfigInterface) *gin.Engine {
	// TODO scope register to scopes-service

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.GetTrustProxy()})

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg, SCOPE_USER_PROTOTYPE)
	// users
	features.LoadUsers(engine, cfg, SCOPE_USER_PROTOTYPE_SELF)

	return engine
}
