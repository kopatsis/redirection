package platform

import (
	"net/http"
	"redir/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
)

func New(db *gorm.DB, ipDB *geoip2.Reader, rdb *redis.Client, httpClient *http.Client) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	router.GET("/", func(c *gin.Context) {
		routes.MainRedirect(c, false)
	})

	router.GET("/:id", routes.Redirect(db, ipDB, rdb, httpClient))

	return router

}
