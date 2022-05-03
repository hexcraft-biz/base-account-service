package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadCommon(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewCommon(cfg)
	e.NoRoute(c.NotFound())

	commonV1 := feature.New(e, "/v1/healthcheck")
	commonV1.GET("/ping", c.Ping())
}
