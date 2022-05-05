package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/base-account-service/middlewares"
	"github.com/hexcraft-biz/feature"
)

func LoadUsers(e *gin.Engine, cfg *config.Config, scopeName string) {
	c := controllers.NewUsers(cfg)

	usersV1 := feature.New(e, "/users/v1")

	usersV1.GET(
		"/users/:id/prototype",
		middlewares.OAuth2PKCE(cfg),
		middlewares.ScopeVerify(cfg, scopeName),
		middlewares.IsSelf(cfg),
		c.Get(),
	)
	usersV1.PUT(
		"/users/:id/prototype/password",
		middlewares.OAuth2PKCE(cfg),
		middlewares.ScopeVerify(cfg, scopeName),
		middlewares.IsSelf(cfg),
		c.UpdatePwd(),
	)
	usersV1.PUT(
		"/users/:id/prototype/status",
		middlewares.OAuth2PKCE(cfg),
		middlewares.ScopeVerify(cfg, scopeName),
		middlewares.IsSelf(cfg),
		c.UpdateStatus(),
	)
	usersV1.DELETE(
		"/users/:id/prototype",
		middlewares.OAuth2PKCE(cfg),
		middlewares.ScopeVerify(cfg, scopeName),
		middlewares.IsSelf(cfg),
		c.Delete(),
	)

}
