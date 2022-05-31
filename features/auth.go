package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/base-account-service/middlewares"
	"github.com/hexcraft-biz/feature"
)

const (
	SCOPE_USER_PROTOTYPE = "user.prototype"
)

func LoadAuth(e *gin.Engine, cfg config.ConfigInterface) {
	c := controllers.NewAuth(cfg)

	authV1 := feature.New(e, "/auth/v1")

	authV1.POST(
		"/login",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.Login(),
	)

	authV1.POST(
		"/signup/confirmation",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.SignUpEmailConfirm(),
	)
	authV1.GET(
		"/signup/tokeninfo",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.SignUpTokenVerify(),
	)
	authV1.POST(
		"/signup",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.SignUp(),
	)

	authV1.POST(
		"/forgetpassword/confirmation",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.ForgetPwdConfirm(),
	)
	authV1.GET(
		"/forgetpassword/tokeninfo",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.ForgetPwdTokenVerify(),
	)
	authV1.PUT(
		"/password",
		middlewares.OAuth2ClientCredentials(cfg),
		middlewares.VerifyScope(cfg, []string{SCOPE_USER_PROTOTYPE}, true),
		c.ChangePassword(),
	)
}
