package middlewares

import (
	"net/http"
	"net/mail"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hexcraft-biz/base-account-service/config"
)

func OAuth2PKCE(cfg config.ConfigInterFace) gin.HandlerFunc {
	/*
		X-{prefix}-Authenticated-User-Email
		X-{prefix}-Authenticated-User-Id
		X-{prefix}-Client-Id
		X-{prefix}-Client-Scope
	*/
	return func(ctx *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()

		authUserEmail := ctx.Request.Header.Get("X-" + prefix + "-Authenticated-User-Email")
		authUserId := ctx.Request.Header.Get("X-" + prefix + "-Authenticated-User-Id")
		clientId := ctx.Request.Header.Get("X-" + prefix + "-Client-Id")
		clientScope := ctx.Request.Header.Get("X-" + prefix + "-Client-Scope")

		if authUserEmail != "" {
			if _, err := mail.ParseAddress(authUserEmail); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
				return
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if authUserId != "" {
			if _, err := uuid.Parse(authUserId); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
				return
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if clientId == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if clientScope == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
	}
}

func OAuth2ClientCredentials(cfg config.ConfigInterFace) gin.HandlerFunc {
	/*
		X-{prefix}-Client-Id
		X-{prefix}-Client-Scope
	*/
	return func(ctx *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()

		clientId := ctx.Request.Header.Get("X-" + prefix + "-Client-Id")
		clientScope := ctx.Request.Header.Get("X-" + prefix + "-Client-Scope")

		if clientId == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if clientScope == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}
	}
}

func IsSelf(cfg config.ConfigInterFace) gin.HandlerFunc {
	/*
		X-{prefix}-Authenticated-User-Id
	*/
	return func(ctx *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()

		authUserId := ctx.Request.Header.Get("X-" + prefix + "-Authenticated-User-Id")

		if authUserId != "" {
			if _, err := uuid.Parse(authUserId); err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
				return
			}
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		userId := ctx.Param("id")
		if authUserId != userId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			return
		}
	}
}

func ScopeVerify(cfg config.ConfigInterFace, scopeName string) gin.HandlerFunc {
	/*
		X-{prefix}-Client-Scope
	*/
	return func(ctx *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()

		clientScope := ctx.Request.Header.Get("X-" + prefix + "-Client-Scope")

		if clientScope == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		} else {
			clientScopes := strings.Split(clientScope, " ")
			sort.Strings(clientScopes)
			if contains(clientScopes, scopeName) == false {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
				return
			}
		}
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
