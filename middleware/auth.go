package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 用户登录，登录后才可以做写操作
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionID := session.Get("id")
		if sessionID == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "unathorized",
			})
		}
	}
}
