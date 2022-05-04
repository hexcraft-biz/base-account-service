package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadUsers(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewUsers(cfg)

	usersV1 := feature.New(e, "/v1/users")
	// TODO middleware check scope & header

	// /users/v1/users/{$userID}/prototype
	usersV1.GET("/me", c.Get())
	usersV1.PATCH("/me", c.Update())
	usersV1.DELETE("/me", c.Delete())
}
