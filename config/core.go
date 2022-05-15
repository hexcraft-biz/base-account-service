package config

import "github.com/jmoiron/sqlx"

type ConfigInterface interface {
	GetDB() *sqlx.DB
	GetJWTSecret() []byte
	GetTrustProxy() string
	GetSMTPHost() string
	GetSMTPPort() string
	GetSMTPUsername() string
	GetSMTPPassword() string
	GetOAuth2HeaderPrefix() string
}
