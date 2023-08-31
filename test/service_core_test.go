package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hexcraft-biz/base-account-service/service"
	"github.com/hexcraft-biz/env"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	cfg, err := Load()
	MustNot(err)
	cfg.DBOpen(false)

	appCfg := &AppConfig{
		DB: cfg.DB,
	}

	router := service.New(appCfg)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck/v1/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	result := struct {
		Message string
	}{}
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		MustNot(err)
	}
	assert.Contains(t, result.Message, "OK")

	/*
		req, _ := http.NewRequest("GET", "/users/v1/users/bd3d395b-7f04-42c7-84f9-9d6effce8f6d/prototype", nil)
		prefix := appCfg.GetOAuth2HeaderPrefix()
		req.Header.Set("X-"+prefix+"-Authenticated-User-Email", "oo7680485@gmail.com")
		req.Header.Set("X-"+prefix+"-Authenticated-User-Id", "bd3d395b-7f04-42c7-84f9-9d6effce8f6d")
		req.Header.Set("X-"+prefix+"-Client-Id", "client_id")
		req.Header.Set("X-"+prefix+"-Client-Scope", "user.prototype")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		result := struct {
			Message string
		}{}
		if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
			MustNot(err)
		}
		assert.Contains(t, result.Message, "OK")
	*/
}

func MustNot(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// ================================================================
// AppConfig implement ConfigInterface
// ================================================================
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

// ================================================================
// Env
// ================================================================
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

// ================================================================
//
// ================================================================
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
