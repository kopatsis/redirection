package routes

import (
	"context"
	"net/http"
	"redir/convert"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
)

func Redirect(db *gorm.DB, ipDB *geoip2.Reader, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		var realURL string

		id, err := convert.FromSixFour(param)
		if err != nil {
			MainRedirect(c)
			return
		}

		cachedURL, err := rdb.Get(context.Background(), param).Result()
		if err == nil {
			realURL = cachedURL
		} else {
			realURL, err = GetRealURL(db, id)
			if err != nil {
				MainRedirect(c)
				return
			}
		}

		go func() {
			click := RequestClickCreate(c, ipDB, id, realURL)
			db.Create(click)
		}()

		c.Redirect(http.StatusSeeOther, realURL)

	}
}
