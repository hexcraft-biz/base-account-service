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
		DB: cfg.DB,
	}

	service.New(appCfg).Run(":" + cfg.Env.AppPort)
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
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

func (ac *AppConfig) GetSignUpEmailSubject() string {
	return os.Getenv("SIGNUP_EMAIL_SUBJECT")
}

func (ac *AppConfig) GetSignUpEmailContent() string {
	return os.Getenv("SIGNUP_EMAIL_CONTENT")
}

func (ac *AppConfig) GetSignUpEmailLinkText() string {
	return os.Getenv("SIGNUP_EMAIL_LINK_TEXT")
}

func (ac *AppConfig) GetForgetPwdEmailSubject() string {
	return os.Getenv("FORGET_PWD_EMAIL_SUBJECT")
}

func (ac *AppConfig) GetForgetPwdEmailContent() string {
	return os.Getenv("FORGET_PWD_EMAIL_CONTENT")
}

func (ac *AppConfig) GetForgetPwdEmailLinkText() string {
	return os.Getenv("FORGET_PWD_LINK_TEXT")
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
