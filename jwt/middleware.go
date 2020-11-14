package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthConfig struct {
	ContextKey  string
	TokenLookup string
	AuthScheme  string
	SigningKey  string
}

var (
	DefaultAuthConfig = AuthConfig{
		ContextKey:  "user",
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
	}
)

type HeaderExtractor struct {
	header     string
	authScheme string
}

func (e HeaderExtractor) ExtractToken(req *http.Request) (string, error) {
	auth := req.Header.Get(e.header)
	l := len(e.authScheme)
	if len(auth) > l+1 && auth[:l] == e.authScheme {
		return auth[l+1:], nil
	}
	return "", errors.New("invalid or expired jwt")
}

func newHeaderExtractor(config AuthConfig) HeaderExtractor {
	parts := strings.Split(config.TokenLookup, ":")
	return HeaderExtractor{
		header:     parts[1],
		authScheme: config.AuthScheme,
	}
}

func AuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.ContextKey, 0)

		extractor := newHeaderExtractor(config)
		token := new(jwt.Token)
		token, err := request.ParseFromRequest(c.Request, extractor, func(token *jwt.Token) (interface{}, error) {
			signingKey := []byte(config.SigningKey)
			return signingKey, nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set(config.ContextKey, token)
		}
	}
}

func HeaderAuthConfig() AuthConfig {
	config := DefaultAuthConfig
	config.SigningKey = getSigningKey()
	return config
}

func HeadAuthHandler() gin.HandlerFunc {
	return AuthMiddleware(HeaderAuthConfig())
}
