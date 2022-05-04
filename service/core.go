package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

type Service struct {
	*gin.Engine
}

func New(cfg config.ConfigInterFace) (*Service, error) {
	// TODO scope design
	// TODO scope register

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.GetTrustProxy()})

	fmt.Println(cfg, cfg.GetDB())

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg)
	// users
	features.LoadUsers(engine, cfg)

	return &Service{Engine: engine}, nil
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
