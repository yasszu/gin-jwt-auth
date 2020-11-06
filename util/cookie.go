package util

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	authKey    = "Authorization"
	expireHour = 24 * 120
)

type CookieStore struct {
	Key        string
	Value      string
	ExpireTime time.Duration
}

func (s CookieStore) Write(c *gin.Context) {
	c.SetCookie(s.Key, s.Value, expireHour, "/", "localhost", false, true)
}

func (s CookieStore) Delete(c *gin.Context) {
	c.SetCookie(s.Key, s.Value, 0, "/", "localhost", false, true)
}

func SaveAuthorizationCookie(token string, c *gin.Context) {
	cookie := CookieStore{Key: authKey, Value: token, ExpireTime: time.Hour * expireHour}
	cookie.Write(c)
}

func DeleteAuthorizationCookie(c *gin.Context) {
	CookieStore{Key: authKey}.Delete(c)
}
