package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func JwtMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}
		if token == "" {
			responseToken(c, "未登录")
			c.Abort()
			return
		}

		parseToken, err := ParseToken(token)
		if err != nil {
			log.Println(err)
			responseToken(c, "token解析失败")
			c.Abort()
			return
		}

		if parseToken.ExpiresAt < time.Now().Unix() {
			responseToken(c, "token过期")
			c.Abort()
			return
		}

		c.Next()
	}
}

func responseToken(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": -1,
		"status_msg":  msg,
	})
}
