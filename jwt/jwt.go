package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims are custom claims extending default ones.
type CustomClaims struct {
	Email     string `json:"email"`
	AccountID uint   `json:"account_id"`
	jwt.StandardClaims
}

const (
	expireHour = 24 * 121
)

func Sign(email string, id uint, secret string) (string, error) {
	expiredAt := time.Now().Add(time.Hour * expireHour).Unix()
	claims := &CustomClaims{
		Email:          email,
		AccountID:      id,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiredAt},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func BindUser(c *gin.Context) *CustomClaims {
	user, _ := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	accountId := claims["account_id"].(float64)
	exp := claims["exp"].(float64)
	return &CustomClaims{
		Email:          email,
		AccountID:      uint(accountId),
		StandardClaims: jwt.StandardClaims{ExpiresAt: int64(exp)},
	}
}

type jwtHeaderExtractor struct {
	header     string
	authScheme string
}

func (e jwtHeaderExtractor) ExtractToken(req *http.Request) (string, error) {
	auth := req.Header.Get(e.header)
	l := len(e.authScheme)
	if len(auth) > l+1 && auth[:l] == e.authScheme {
		return auth[l+1:], nil
	}
	return "", errors.New("invalid or expired jwt")
}

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

func NewAuthConfig(signingKey string) AuthConfig {
	config := DefaultAuthConfig
	config.SigningKey = signingKey
	return config
}

func AuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(config.ContextKey, 0)

		parts := strings.Split(config.TokenLookup, ":")
		extractor := jwtHeaderExtractor{parts[1], config.AuthScheme}

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
