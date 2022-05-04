package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

type Service struct {
	engine *gin.Engine
	config *config.Config
}

func New() (*Service, error) {
	cfg, err := config.Load()
	MustNot(err)

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.TrustProxy})

	MustNot(cfg.DBOpen(false))

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg)
	// users
	features.LoadUsers(engine, cfg)

	return &Service{engine: engine, config: cfg}, nil
}

func (s *Service) GetEngine() *gin.Engine {
	return s.engine
}

func (s *Service) GetConfig() *config.Config {
	return s.config
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
