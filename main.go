package main

import (
	"os"

	serviceConfig "github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/service"
	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := Load()
	MustNot(err)
	cfg.DBOpen(false)

	sc := &serviceConfig.Config{
		DB: cfg.DB,
		Env: &serviceConfig.Env{
			JWTSecret:    []byte(os.Getenv("JWT_SECRET")),
			TrustProxy:   cfg.TrustProxy,
			SMTPHost:     os.Getenv("SMTP_HOST"),
			SMTPPort:     os.Getenv("SMTP_PORT"),
			SMTPUsername: os.Getenv("SMTP_USERNAME"),
			SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		},
	}

	service.New(sc).Run(":" + cfg.Env.AppPort)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//================================================================
// Env
//================================================================
type Env struct {
	*env.Prototype
}

func FetchEnv() (*Env, error) {
	if e, err := env.Fetch(); err != nil {
		return nil, err
	} else {
		return &Env{
			Prototype: e,
		}, nil
	}
}

//================================================================
//
//================================================================
type Config struct {
	*Env
	DB *sqlx.DB
}

func Load() (*Config, error) {
	e, err := FetchEnv()
	if err != nil {
		return nil, err
	}

	return &Config{Env: e}, nil
}

func (cfg *Config) DBOpen(init bool) error {
	var err error

	cfg.DBClose()
	cfg.DB, err = cfg.MysqlConnectWithMode(init)

	return err
}

func (cfg *Config) DBClose() {
	if cfg.DB != nil {
		cfg.DB.Close()
	}
}
