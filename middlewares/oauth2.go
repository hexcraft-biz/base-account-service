package middlewares

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hexcraft-biz/base-account-service/config"
	"github.com/hexcraft-biz/base-account-service/models"
)

const (
	ScopeDelimiter = " "
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

		if authUserEmail == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		if authUserId == "" {
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
		authUserEmail := ctx.Request.Header.Get("X-" + prefix + "-Authenticated-User-Email")

		// ID
		userId := ctx.Param("id")
		if authUserId != userId {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			return
		}

		// Email
		if entityRes, err := models.NewUsersTableEngine(cfg.GetDB()).GetByID(authUserId); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		} else {
			if entityRes == nil || authUserEmail != entityRes.Identity {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
				return
			}
		}
	}
}

func ScopeVerify(cfg config.ConfigInterFace, resourceScopes []string, isExact bool) gin.HandlerFunc {
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

			if scopeIntersect(clientScopes, resourceScopes, isExact) == false {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
				return
			}
		}
	}
}

func scopeIntersect(clientScopes, resourceScopes []string, isExact bool) bool {
	hashTable := make(map[string]int)

	newClientScopes := removeDuplicateStr(clientScopes)
	newResourceScopes := removeDuplicateStr(resourceScopes)

	for i := range newClientScopes {
		hashTable[newClientScopes[i]]++
	}

	for i := range newResourceScopes {
		hashTable[newResourceScopes[i]]++
	}

	matchCount := 0
	for _, s := range newResourceScopes {
		if hashTable[s] >= 2 {
			matchCount++
		}
	}

	if isExact == true {
		return matchCount == len(resourceScopes)
	} else {
		return matchCount >= 1
	}
}

func removeDuplicateStr(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//================================================================
//
//================================================================
func VerifyScope(cfg config.ConfigInterFace, allows []string) gin.HandlerFunc {
	/*
		X-{prefix}-Client-Scope
	*/
	return func(ctx *gin.Context) {
		clientScopes := strings.Split(ctx.Request.Header.Get("X-"+cfg.GetOAuth2HeaderPrefix()+"-Client-Scope"), ScopeDelimiter)
		if !inAllows(allows, clientScopes) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
		}
	}
}

func inAllows(clientScopes, allows []string) bool {
	sort.Strings(allows)
	l := len(allows)
	for i := range clientScopes {
		if sort.SearchStrings(allows, clientScopes[i]) < l {
			return true
		}
	}
	return false
}
