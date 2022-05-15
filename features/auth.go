package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/base-account-service/middlewares"
	"github.com/hexcraft-biz/feature"
)

func LoadAuth(e *gin.Engine, cfg config.ConfigInterface, scopeName string) {
	c := controllers.NewAuth(cfg)

	authV1 := feature.New(e, "/auth/v1")

	authV1.POST(
		"/login",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.Login(),
	)

	authV1.POST(
		"/signup/confirmation",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.SignUpEmailConfirm(),
	)
	authV1.GET(
		"/signup/tokeninfo",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.SignUpTokenVerify(),
	)
	authV1.POST(
		"/signup",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.SignUp(),
	)

	authV1.POST(
		"/forgetpassword/confirmation",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.ForgetPwdConfirm(),
	)
	authV1.GET(
		"/forgetpassword/tokeninfo",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.ForgetPwdTokenVerify(),
	)
	authV1.PUT(
		"/password",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.ScopeVerify(cfg, []string{scopeName}, true),
		c.ChangePassword(),
	)
}
