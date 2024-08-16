package platform

import "github.com/gin-gonic/gin"

func New() *gin.Engine {
	router := gin.Default()
	return router

}
