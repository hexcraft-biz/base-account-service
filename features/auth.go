package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadAuth(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewAuth(cfg)

	// TODO middleware check scope

	authV1 := feature.New(e, "/v1/auth")

	authV1.POST("/login", c.Login()) // ok

	authV1.POST("/signup/confirmation", c.SignUpEmailComfirm()) // ok //TODO send email
	authV1.GET("/signup/tokeninfo", c.SignUpTokenVerify())      // ok
	authV1.POST("/signup", c.SignUp())                          // ok

	authV1.POST("/resetpassword/confirmation", c.ResetPwdComfirm()) // ok //TODO send email
	authV1.GET("/resetpassword/tokeninfo", c.ResetPwdTokenVerify()) // ok
	authV1.PUT("/passoword", c.ChangePassword())                    // ok
}
