package config

import "github.com/jmoiron/sqlx"

type Config struct {
	DB  *sqlx.DB
	Env *Env
}

type Env struct {
	JWTSecret          []byte
	TrustProxy         string
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	OAuth2HeaderPrefix string
}
