package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/features"
)

type Service struct {
	Engine *gin.Engine
	Config *config.Config
}

func New() (*Service, error) {
	cfg, err := config.Load()
	MustNot(err)

	engine := gin.Default()
	engine.SetTrustedProxies([]string{cfg.TrustProxy})
	// TODO session cookies issue

	if cfg.AutoCreateDBSchema {
		MustNot(cfg.MysqlDBInit("./sql/"))
	}

	MustNot(cfg.DBOpen(false))
	// defer cfg.DBClose()

	return &Service{Engine: engine, Config: cfg}, nil
}

func (s *Service) Run() {

	// base features
	features.LoadCommon(s.Engine, s.Config)
	// auth
	features.LoadAuth(s.Engine, s.Config)
	// users
	features.LoadUsers(s.Engine, s.Config)

	s.Engine.Run(":" + s.Config.AppPort)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
