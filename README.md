# base-account-service
The base-account-service for building a customer account system.  
You can inherit from base-account-service and extend and develop the account system you need.  
Please remember to use the same database with [accounts-service-backend](https://github.com/hexcraft-biz/accounts-service-backend)  

## Endpoint
### HealthCheck
#### GET /healthcheck/v1/ping
- Params
  - None
- Response
  - 200
	```json
	{
	  "message": "OK"
	}
	```

### UserPrototype
#### GET /users/v1/users/:user_id/prototype
- Required Scope : `user.prototype`
- Params
  - Headers
    - Authorization : Bearer {oauth2_token}
    - X-{OAUTH2_HEADER_PREFIX}-Client-Id : {client_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Client-Scope : {scope} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Id : {user_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Email : {user_email} (will generate from api-proxy by oauth2_token)
  - Uri
    - user_id
      - Required : True
      - Type : String
      - Example : "9cfa987b-022d-4461-82c6-f7f12d706163"
- Response
  - 200
	```json
	{
	  "id": "9cfa987b-022d-4461-82c6-f7f12d706163",
	  "identity": "xxx@mail.com",
	  "status": "enabled",
	  "createdAt": "2022-11-01 07:08:34",
	  "updatedAt": "2022-11-01 07:08:34"
	}
	```
  - 400 | 401 | 403 | 404 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### PUT /users/v1/users/:user_id/prototype/password
- Required Scope : `user.prototype`
- Params
  - Headers
    - Authorization : Bearer {oauth2_token}
    - Content-Type : application/json
    - X-{OAUTH2_HEADER_PREFIX}-Client-Id : {client_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Client-Scope : {scope} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Id : {user_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Email : {user_email} (will generate from api-proxy by oauth2_token)
  - Body
    - password
      - Required : True
      - Type : String
      - Example : "IamPassword"
- Response
  - 204
  - 400 | 401 | 403 | 404 | 409 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

#### PUT /users/v1/users/:user_id/prototype/status
- Required Scope : `user.prototype`
- Params
  - Headers
    - Authorization : Bearer {oauth2_token}
    - Content-Type : application/json
    - X-{OAUTH2_HEADER_PREFIX}-Client-Id : {client_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Client-Scope : {scope} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Id : {user_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Email : {user_email} (will generate from api-proxy by oauth2_token)
  - Body
    - status
      - Required : True
      - Allowed : `enabled` `disabled` `suspended`
- Response
  - 200
	```json
	{
	  "id": "9cfa987b-022d-4461-82c6-f7f12d706163",
	  "identity": "xxx@mail.com",
	  "status": "enabled",
	  "createdAt": "2022-11-01 07:08:34",
	  "updatedAt": "2022-11-01 07:08:34"
	}
	```
  - 400 | 401 | 403 | 404 | 500
	```json
	{
	  "message": "Error Message"
	}
	```
#### DELETE /users/v1/users/:user_id/prototype
- Required Scope : `user.prototype`
- Params
  - Headers
    - Authorization : Bearer {oauth2_token}
    - X-{OAUTH2_HEADER_PREFIX}-Client-Id : {client_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Client-Scope : {scope} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Id : {user_id} (will generate from api-proxy by oauth2_token)
    - X-{OAUTH2_HEADER_PREFIX}-Authenticated-User-Email : {user_email} (will generate from api-proxy by oauth2_token)
  - Uri
    - user_id
      - Required : True
      - Type : String
      - Example : "9cfa987b-022d-4461-82c6-f7f12d706163"
- Response
  - 204
  - 400 | 401 | 403 | 404 | 500
	```json
	{
	  "message": "Error Message"
	}
	```

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
