package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/misc"
	"github.com/hexcraft-biz/base-account-service/models"
	"github.com/hexcraft-biz/controller"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_STATUS_ENABLED            = "enabled"
	EMAIL_CONFIRMATION_EXPIRE_MINS = 10
	JWT_TYPE_SIGN_UP               = "signup"
	JWT_TYPE_RESET_PWD             = "resetpwd"
)

type Auth struct {
	*controller.Prototype
	Config config.ConfigInterFace
}

func NewAuth(cfg config.ConfigInterFace) *Auth {
	return &Auth{
		Prototype: controller.New("auth", cfg.GetDB()),
		Config:    cfg,
	}
}

//================================================================
// Auth Login
//================================================================
type genTokenParams struct {
	Identity string `json:"identity" binding:"required,email,min=1,max=128"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params genTokenParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Identity); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				saltedPwd := append([]byte(params.Password), entityRes.Salt...)
				compareErr := bcrypt.CompareHashAndPassword(entityRes.Password, saltedPwd)
				if compareErr != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Password is wrong."})
					return
				}

				if entityRes.Status != USER_STATUS_ENABLED {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "This account is not enabled."})
					return
				}

				if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
					return
				} else {
					c.JSON(http.StatusOK, absRes)
					return
				}
			}
		}
	}
}

//================================================================
// SignUp
//================================================================
type signUpEmailConfirmParams struct {
	Email         string `json:"email" binding:"required,email,min=1,max=128"`
	VerifyPageUrl string `json:"verify_page_url" binding:"required,url"`
}

type signUpEmailConfirmResp struct {
	Token string `json:"token"`
}

func (ctrl *Auth) SignUpEmailConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signUpEmailConfirmParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
			return
		} else if entityRes != nil {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "This Email is already exist."})
			return
		}

		nowTime := time.Now()
		expiresAt := nowTime.Add(EMAIL_CONFIRMATION_EXPIRE_MINS * time.Minute).Unix()
		issuedAt := nowTime.Unix()

		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		tokenString, err := miscJWT.GenToken(jwt.SigningMethodHS512, misc.EmailJwtClaims{
			StandardClaims: jwt.StandardClaims{
				Subject:   params.Email,
				ExpiresAt: expiresAt,
				IssuedAt:  issuedAt,
			},
			Email: params.Email,
			Type:  JWT_TYPE_SIGN_UP,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		email := misc.NewEmail(
			ctrl.Config.GetSMTPHost(),
			ctrl.Config.GetSMTPPort(),
			ctrl.Config.GetSMTPUsername(),
			ctrl.Config.GetSMTPPassword(),
		)
		to := []string{params.Email}
		subject := "Signup Email Confirmation"
		body := `<html><body>This is email confirmation, please follow this <a href="` + params.VerifyPageUrl + tokenString + `">link</a> to complete sign up flow.</body></html>`
		email.SendHTML(to, subject, body)

		c.JSON(http.StatusAccepted, gin.H{"message": http.StatusText(http.StatusAccepted)})
		return
	}
}

type signUpTokenVerifyParams struct {
	Token string `form:"token" binding:"required"`
}

type signUpTokenVerifyResp struct {
	Email string `json:"email"`
}

func (ctrl *Auth) SignUpTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signUpTokenVerifyParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_SIGN_UP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.JSON(http.StatusOK, signUpTokenVerifyResp{
			Email: claims.Email,
		})
		return
	}
}

type signupParams struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signupParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_SIGN_UP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).Insert(claims.Email, params.Password, USER_STATUS_ENABLED); err != nil {
			if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
				c.JSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}
		} else {
			if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": absErr.Error()})
			} else {
				c.JSON(http.StatusOK, absRes)
			}
		}
	}
}

//================================================================
// ResetPassword
//================================================================
type resetPwdConfirmParams struct {
	Email         string `json:"email" binding:"required,email,min=1,max=128"`
	VerifyPageUrl string `json:"verify_page_url" binding:"required,url"`
}

type resetPwdConfirmResp struct {
	Token string `json:"token"`
}

func (ctrl *Auth) ResetPwdConfirm() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params resetPwdConfirmParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
			return
		} else if entityRes == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "This Email is not already exist."})
			return
		}

		nowTime := time.Now()
		expiresAt := nowTime.Add(EMAIL_CONFIRMATION_EXPIRE_MINS * time.Minute).Unix()
		issuedAt := nowTime.Unix()

		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		tokenString, err := miscJWT.GenToken(jwt.SigningMethodHS512, misc.EmailJwtClaims{
			StandardClaims: jwt.StandardClaims{
				Subject:   params.Email,
				ExpiresAt: expiresAt,
				IssuedAt:  issuedAt,
			},
			Email: params.Email,
			Type:  JWT_TYPE_RESET_PWD,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		email := misc.NewEmail(
			ctrl.Config.GetSMTPHost(),
			ctrl.Config.GetSMTPPort(),
			ctrl.Config.GetSMTPUsername(),
			ctrl.Config.GetSMTPPassword(),
		)
		to := []string{params.Email}
		subject := "Reset Password Email Confirmation"
		body := `<html><body>This is email confirmation, please follow this <a href="` + params.VerifyPageUrl + tokenString + `">link</a> to complete reset password flow.</body></html>`
		email.SendHTML(to, subject, body)

		c.JSON(http.StatusAccepted, gin.H{"message": http.StatusText(http.StatusAccepted)})
		return
	}
}

type resetPwdTokenVerifyParams struct {
	Token string `form:"token" binding:"required"`
}

type resetPwdTokenVerifyResp struct {
	Email string `json:"email"`
}

func (ctrl *Auth) ResetPwdTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params resetPwdTokenVerifyParams
		if err := c.ShouldBindQuery(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_RESET_PWD {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		c.JSON(http.StatusOK, resetPwdTokenVerifyResp{
			Email: claims.Email,
		})
		return
	}
}

type resetPwdParams struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params resetPwdParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(ctrl.Config.GetJWTSecret())
		token, err := miscJWT.Parse(params.Token, &claims)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_RESET_PWD {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		usersEngine := models.NewUsersTableEngine(ctrl.DB)

		if entityRes, err := usersEngine.GetByIdentity(claims.Email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
				return
			} else {
				// TODO next version about password log
				saltedPwd := append([]byte(params.Password), entityRes.Salt...)
				compareErr := bcrypt.CompareHashAndPassword(entityRes.Password, saltedPwd)
				if compareErr == nil {
					c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
					return
				}

				if rowsAffected, err := usersEngine.ResetPwd(entityRes.ID, params.Password, entityRes.Salt); err != nil {
					if myErr, ok := err.(*mysql.MySQLError); ok && myErr.Number == 1062 {
						c.JSON(http.StatusConflict, gin.H{"message": http.StatusText(http.StatusConflict)})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					}
				} else {
					if rowsAffected == 0 {
						c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
						return
					} else {
						c.AbortWithStatusJSON(http.StatusNoContent, gin.H{"message": http.StatusText(http.StatusNoContent)})
						return
					}
				}

			}
		}
	}
}
