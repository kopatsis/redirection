package routes

import "github.com/gin-gonic/gin"

func Redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.ClientIP()
	}
}
