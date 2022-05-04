package config

import (
	"os"

	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
)

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
