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

func OAuth2PKCE(cfg config.ConfigInterface) gin.HandlerFunc {
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

func OAuth2ClientCredentials(cfg config.ConfigInterface) gin.HandlerFunc {
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

func IsSelf(cfg config.ConfigInterface, selfScope string, allowScopes []string) gin.HandlerFunc {
	/*
		X-{prefix}-Authenticated-User-Id
	*/
	return func(c *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()
		authUserId := c.Request.Header.Get("X-" + prefix + "-Authenticated-User-Id")
		authUserEmail := c.Request.Header.Get("X-" + prefix + "-Authenticated-User-Email")
		clientScope := strings.Split(c.Request.Header.Get("X-"+prefix+"-Client-Scope"), ScopeDelimiter)
		reqUserID := c.Param("id")

		if HasScope(selfScope, clientScope) {
			if authUserId != reqUserID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			} else if entityRes, err := models.NewUsersTableEngine(cfg.GetDB()).GetByID(authUserId); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else if entityRes == nil || authUserEmail != entityRes.Identity {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			}
		} else if !InAllows(allowScopes, clientScope) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
		}
	}
}

func VerifyScope(cfg config.ConfigInterface, resourceScopes []string, isExact bool) gin.HandlerFunc {
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
			clientScopes := strings.Split(clientScope, ScopeDelimiter)

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
type Account interface {
	GetIdentity() string
}

type UserAccounts interface {
	GetMwInterfaceByID(userID string) (Account, error)
}

func IsSelfRequest(cfg config.ConfigInterface, mei UserAccounts, selfScope string, allowScopes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		prefix := cfg.GetOAuth2HeaderPrefix()
		authUserID := c.Request.Header.Get("X-" + prefix + "-Authenticated-User-Id")
		authUserEmail := c.Request.Header.Get("X-" + prefix + "-Authenticated-User-Email")
		clientScope := strings.Split(c.Request.Header.Get("X-"+prefix+"-Client-Scope"), ScopeDelimiter)
		userID := c.Param("id")

		if HasScope(selfScope, clientScope) {
			if authUserID != userID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			} else if row, err := mei.GetMwInterfaceByID(userID); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else if row == nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			} else if authUserEmail != row.GetIdentity() {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
			} else {
				c.Set("user", row)
			}
		} else if InAllows(allowScopes, clientScope) {
			if row, err := mei.GetMwInterfaceByID(userID); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			} else if row == nil {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": http.StatusText(http.StatusNotFound)})
			} else {
				c.Set("user", row)
			}
		}
	}
}

func HasScope(s string, scopes []string) bool {
	for i := range scopes {
		if s == scopes[i] {
			return true
		}
	}
	return false
}

//================================================================
//
//================================================================
func VerifyScopeWithHeaderAffix(headerAffix string, allows []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientScopes := strings.Split(c.Request.Header.Get("X-"+headerAffix+"-Client-Scope"), ScopeDelimiter)
		if !InAllows(allows, clientScopes) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": http.StatusText(http.StatusForbidden)})
		}
	}
}

func InAllows(allows, scopes []string) bool {
	sort.Strings(allows)
	l := len(allows)
	for i := range scopes {
		x := sort.SearchStrings(allows, scopes[i])
		if (x < l) && (allows[x] == scopes[i]) {
			return true
		}
	}
	return false
}
