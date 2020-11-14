package jwt

import (
	"gin-jwt-auth/model"
	"github.com/gin-gonic/gin"
	"os"
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
