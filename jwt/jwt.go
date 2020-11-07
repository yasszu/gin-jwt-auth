package jwt

import (
	"errors"
	"gin-jwt-auth/model"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims are custom claims extending default ones.
type CustomClaims struct {
	AccountID uint `json:"account_id"`
	jwt.StandardClaims
}

const (
	expireHour = 24 * 121
)

func getSigningKey() string {
	os.Setenv("JWT_SECRET", "b5a636fc-bd01-41b1-9780-7bbd906fa4c0")
	return os.Getenv("JWT_SECRET")
}

func Sign(account *model.Account) (*model.AccessToken, error) {
	expiredAt := time.Now().Add(time.Hour * expireHour)
	claims := &CustomClaims{
		AccountID:      account.ID,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expiredAt.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(getSigningKey()))
	if err != nil {
		return nil, err
	}
	accessToken := &model.AccessToken{
		AccountID: account.ID,
		Token:     signedString,
		ExpiresAt: expiredAt,
	}
	return accessToken, err
}

func BindUser(c *gin.Context) *CustomClaims {
	user, _ := c.Get("user")
	token := user.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(float64)
	exp := claims["exp"].(float64)
	return &CustomClaims{
		AccountID:      uint(accountId),
		StandardClaims: jwt.StandardClaims{ExpiresAt: int64(exp)},
	}
}

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

func HeaderAuthConfig() AuthConfig {
	config := DefaultAuthConfig
	config.SigningKey = getSigningKey()
	return config
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
