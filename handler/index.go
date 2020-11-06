package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Index Handler
func Index(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
