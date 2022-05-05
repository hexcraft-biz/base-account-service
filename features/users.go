package features

import (
	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/controllers"
	"github.com/hexcraft-biz/feature"
)

func LoadUsers(e *gin.Engine, cfg *config.Config) {
	c := controllers.NewUsers(cfg)

	usersV1 := feature.New(e, "/users/v1")
	// TODO middleware check scope & header
	/*
		if scope no match, reject 403
		scope : user.prototype.self

		if header not contain these, reject 401
		X-Kmk-Authenticated-User-Email
		X-Kmk-Authenticated-User-Id
		X-Kmk-Client-Id
		X-Kmk-Client-Scope

		if uri.ID not match X-Kmk-Authenticated-User-Id, reject 403
	*/

	usersV1.GET("/users/:id/prototype", c.Get())                 // OK
	usersV1.PUT("/users/:id/prototype/password", c.UpdatePwd())  // OK
	usersV1.PUT("/users/:id/prototype/status", c.UpdateStatus()) // OK
	usersV1.DELETE("/users/:id/prototype", c.Delete())           // OK

}
