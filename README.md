# base-account-service
The base-account-service for building a customer account system.  
You can inherit from base-account-service and extend and develop the account system you need.  
Please remember to use the same database with [accounts-service-backend](https://github.com/hexcraft-biz/accounts-service-backend)  

## Quick start

```sh
# assume the following codes in main.go file
$ cat main.go
```

```go
package main

import (
	"github.com/hexcraft-biz/base-account-service/service"
)

func main() {
	appCfg := &AppConfig{
		DB: cfg.DB,
	}

	engine := service.New(appCfg)

	// TODO
	// Do something...
	// Like using gin.Engine

	engine.Run(":" + appPort)
}

//================================================================
// AppConfig implement ConfigInterface
//================================================================
type AppConfig struct {
	DB *sqlx.DB
}

func (ac *AppConfig) GetDB() *sqlx.DB {
	return ac.DB
}

func (ac *AppConfig) GetTrustProxy() string {
	return os.Getenv("TRUST_PROXY")
}

func (ac *AppConfig) GetOAuth2HeaderPrefix() string {
	return os.Getenv("OAUTH2_HEADER_PREFIX")
}

```
