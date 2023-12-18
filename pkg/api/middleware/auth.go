package middleware

import (
	"HeadZone/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AdminAuthMiddleware(c *gin.Context) {

	accessToken := c.Request.Header.Get("Authorization")

	cfg, _ := config.LoadConfig()

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.ACCESS_KEY_ADMIN), nil
	})

	if err != nil {
		// The access token is invalid.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
		c.Abort()
		return
		// c.AbortWithStatus(401)
		// return
	}

	c.Next()
}
