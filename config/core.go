package config

import "github.com/jmoiron/sqlx"

type ConfigInterface interface {
	GetDB() *sqlx.DB
	GetTrustProxy() string
	GetOAuth2HeaderPrefix() string
}
