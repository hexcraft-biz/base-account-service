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

	//if cfg.AutoCreateDBSchema {
	//	MustNot(cfg.MysqlDBInit("./sql/"))
	//}

	//MustNot(cfg.DBOpen(false))
	// defer cfg.DBClose()

	// base features
	features.LoadCommon(engine, cfg)
	// auth
	features.LoadAuth(engine, cfg)
	// users
	features.LoadUsers(engine, cfg)

	return &Service{Engine: engine, Config: cfg}, nil
}

func (s *Service) GetEngine() *gin.Engine {
	return s.Engine
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}
