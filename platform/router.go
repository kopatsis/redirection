package platform

import (
	"redir/routes"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
)

func New(db *gorm.DB, ipDB *geoip2.Reader) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/", routes.MainRedirect)

	router.GET("/:id", routes.Redirect(db, ipDB))

	return router

}
