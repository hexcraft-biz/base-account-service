package main

import (
	"os"

	"github.com/hexcraft-biz/base-account-service/service"
	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := Load()
	MustNot(err)
	cfg.DBOpen(false)

	appCfg := &AppConfig{
		DB:         cfg.DB,
		JWTSecret:  cfg.Env.JwtSecret,
		TrustProxy: cfg.TrustProxy,
	}

	service.New(appCfg).Run(":" + cfg.Env.AppPort)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//================================================================
// AppConfig implement ConfigInterFace
//================================================================
type AppConfig struct {
	DB         *sqlx.DB
	JWTSecret  []byte
	TrustProxy string
}

func (ac *AppConfig) GetDB() *sqlx.DB {
	return ac.DB
}

func (ac *AppConfig) GetJWTSecret() []byte {
	return ac.JWTSecret
}

func (ac *AppConfig) GetTrustProxy() string {
	return ac.TrustProxy
}

//================================================================
// Env
//================================================================
type Env struct {
	*env.Prototype
	JwtSecret []byte
}

func FetchEnv() (*Env, error) {
	if e, err := env.Fetch(); err != nil {
		return nil, err
	} else {
		return &Env{
			Prototype: e,
			JwtSecret: []byte(os.Getenv("JWT_SECRET")),
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
