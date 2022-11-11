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
	GetSMTPSender() string
	GetOAuth2HeaderPrefix() string
	GetSignUpEmailSubject() string
	GetSignUpEmailContent() string
	GetSignUpEmailLinkText() string
	GetForgetPwdEmailSubject() string
	GetForgetPwdEmailContent() string
	GetForgetPwdEmailLinkText() string
}
