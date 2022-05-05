package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

func New(cfg *config.Config) *gin.Engine {
	// TODO scope design
	// TODO scope register

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.Env.TrustProxy})

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg)
	// users
	features.LoadUsers(engine, cfg)

	return engine
}
