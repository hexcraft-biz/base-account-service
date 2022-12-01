package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

func New(cfg config.ConfigInterface) *gin.Engine {
	// TODO scope register to scopes-service

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.GetTrustProxy()})

	// base features
	features.LoadCommon(engine, cfg)
	// users
	features.LoadUsers(engine, cfg)

	return engine
}
