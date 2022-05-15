# Base Account Service

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

func (ac *AppConfig) GetJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func (ac *AppConfig) GetTrustProxy() string {
	return os.Getenv("TRUST_PROXY")
}

func (ac *AppConfig) GetSMTPHost() string {
	return os.Getenv("SMTP_HOST")
}

func (ac *AppConfig) GetSMTPPort() string {
	return os.Getenv("SMTP_PORT")
}

func (ac *AppConfig) GetSMTPUsername() string {
	return os.Getenv("SMTP_USERNAME")
}

func (ac *AppConfig) GetSMTPPassword() string {
	return os.Getenv("SMTP_PASSWORD")
}

func (ac *AppConfig) GetOAuth2HeaderPrefix() string {
	return os.Getenv("OAUTH2_HEADER_PREFIX")
}

```
