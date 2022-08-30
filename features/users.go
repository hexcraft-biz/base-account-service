package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/base-account-service/middlewares"
	"github.com/hexcraft-biz/feature"
)

const (
	SCOPE_USER_PROTOTYPE_SELF = "user.prototype.self"
)

func LoadUsers(e *gin.Engine, cfg config.ConfigInterface) {
	c := controllers.NewUsers(cfg)

	usersV1 := feature.New(e, "/users/v1")

	usersV1.GET(
		"/users/:id/prototype",
		middlewares.OAuth2PKCE(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE_SELF}, true),
		middlewares.IsSelf(cfg, SCOPE_USER_PROTOTYPE_SELF, []string{}),
		c.Get(),
	)
	usersV1.PUT(
		"/users/:id/prototype/password",
		middlewares.OAuth2PKCE(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE_SELF}, true),
		middlewares.IsSelf(cfg, SCOPE_USER_PROTOTYPE_SELF, []string{}),
		c.UpdatePwd(),
	)
	usersV1.PUT(
		"/users/:id/prototype/status",
		middlewares.OAuth2PKCE(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE_SELF}, true),
		middlewares.IsSelf(cfg, SCOPE_USER_PROTOTYPE_SELF, []string{}),
		c.UpdateStatus(),
	)
}
