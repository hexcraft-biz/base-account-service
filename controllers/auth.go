package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/misc"
	"github.com/hexcraft-biz/base-account-service/models"
	"github.com/hexcraft-biz/controller"
)

const (
	NEW_REGISTERED_USER_STATUS     = "suspended"
	EMAIL_CONFIRMATION_EXPIRE_MINS = 10
	JWT_TYPE_SIGN_UP               = "signup"
)

// TODO move to env
var jwtKey = []byte("FDr1VjVQiSiybYJrQZNt8Vfd7bFEsKP6vNX1brOSiWl0mAIVCxJiR4/T3zpAlBKc2/9Lw2ac4IwMElGZkssfj3dqwa7CQC7IIB+nVxiM1c9yfowAZw4WQJ86RCUTXaXvRX8JoNYlgXcRrK3BK0E/fKCOY1+izInW3abf0jEeN40HJLkXG6MZnYdhzLnPgLL/TnIFTTAbbItxqWBtkz6FkZTG+dkDSXN7xNUxlg==")

type Auth struct {
	*controller.Prototype
}

func NewAuth(cfg *config.Config) *Auth {
	return &Auth{
		Prototype: controller.New("auth", cfg.DB),
	}
}

//================================================================
// Auth Token
//================================================================
type genTokenParams struct {
	Identity string `json:"identity" binding:"required,email,min=1,max=128"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) GenerateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params genTokenParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).GetByIdentity(params.Identity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": err.Error()})
		} else {
			if entityRes == nil {
				fmt.Println(entityRes)
			}
			if absRes, absErr := entityRes.GetAbsUser(); absErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError), "results": absErr.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK), "results": absRes})
			}
		}
	}
}

//================================================================
// SignUp
//================================================================
type signUpEmailComfirmParams struct {
	Email string `json:"email" binding:"required,email,min=1,max=128"`
}

type signUpEmailComfirmResp struct {
	Token string `json:"token"`
}

func (ctrl *Auth) SignUpEmailComfirm() gin.HandlerFunc {
	return func(c *gin.Context) {

		var params signUpEmailComfirmParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// TODO check email exist or not.

		nowTime := time.Now()
		expiresAt := nowTime.Add(EMAIL_CONFIRMATION_EXPIRE_MINS * time.Minute).Unix()
		issuedAt := nowTime.Unix()

		miscJWT := misc.NewJWT(jwtKey)
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

		// TODO Send Comfirmation Email
		c.JSON(http.StatusOK, signUpEmailComfirmResp{
			Token: tokenString,
		})
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
		miscJWT := misc.NewJWT(jwtKey)
		token, err := miscJWT.Parse(params.Token, &claims, jwtKey)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
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

type SignupParams struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=5,max=128"`
}

func (ctrl *Auth) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		var signupParams SignupParams
		if err := c.ShouldBindJSON(&signupParams); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var claims misc.EmailJwtClaims
		miscJWT := misc.NewJWT(jwtKey)
		token, err := miscJWT.Parse(signupParams.Token, &claims, jwtKey)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		if !token.Valid || claims.Type != JWT_TYPE_SIGN_UP {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if entityRes, err := models.NewUsersTableEngine(ctrl.DB).Insert(claims.Email, signupParams.Password, NEW_REGISTERED_USER_STATUS); err != nil {
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
func (ctrl *Auth) ResetPwdComfirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}

func (ctrl *Auth) ResetPwdTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}

func (ctrl *Auth) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO get email from JWT
		c.JSON(http.StatusOK, gin.H{"message": http.StatusText(http.StatusOK)})
	}
}
