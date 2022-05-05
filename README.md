# Base Account Service

## Quick start

```sh
# assume the following codes in main.go file
$ cat example.go
```

```go
package main

import (
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/service"
)

func main() {
	sc := &serviceConfig.Config{
		DB: cfg.DB,
		Env: &serviceConfig.Env{
			JWTSecret:          []byte(os.Getenv("JWT_SECRET")),
			TrustProxy:         os.Getenv("TRUST_PROXY"),
			SMTPHost:           os.Getenv("SMTP_HOST"),
			SMTPPort:           os.Getenv("SMTP_PORT"),
			SMTPUsername:       os.Getenv("SMTP_USERNAME"),
			SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
			OAuth2HeaderPrefix: os.Getenv("OAUTH2_HEADER_PREFIX"),
		},
	}

	engine := service.New(sc)

	// TODO
	// Do something...
	// Like using gin.Engine

	engine.Run(":" + cfg.Env.AppPort)
}
```
