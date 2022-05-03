package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadAuth(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewAuth(cfg)

	authV1 := feature.New(e, "/v1/auth")

	authV1.POST("/token", c.GenerateToken())

	authV1.POST("/signup/confirmation", c.SignUpEmailComfirm()) // ok
	authV1.GET("/signup/tokeninfo", c.SignUpTokenVerify())      // ok
	authV1.POST("/signup", c.SignUp())                          // ok

	authV1.POST("/resetpassword/confirmation", c.ResetPwdComfirm())
	authV1.GET("/resetpassword/tokeninfo", c.ResetPwdTokenVerify())
	authV1.PUT("/passoword", c.ChangePassword())
}
