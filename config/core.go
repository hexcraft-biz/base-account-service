package config

import "github.com/jmoiron/sqlx"

type ConfigInterFace interface {
	GetDB() *sqlx.DB
	GetJWTSecret() []byte
	GetTrustProxy() string
}
