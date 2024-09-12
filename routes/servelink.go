package routes

import (
	"context"
	"net/http"
	"redir/convert"
	"redir/datatypes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
)

func Redirect(db *gorm.DB, ipDB *geoip2.Reader, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		var realURL string

		var id int
		custom := len(param) > 6

		if !custom {
			var err error
			id, err = convert.FromSixFour(param)
			if err != nil {
				MainRedirect(c, true)
				return
			}
		}

		cachedURL, err := rdb.Get(context.Background(), param).Result()
		if (len(cachedURL) >= 3 && cachedURL[:3] == ":e:") || cachedURL == ":a:" {
			MainRedirect(c, true)
			return
		} else if err == nil {
			realURL = cachedURL
		} else {
			if custom {
				realURL, err = GetRealURLByCustom(db, param)
			} else {
				realURL, err = GetRealURL(db, id)
			}

			if err != nil {
				MainRedirect(c, true)
				return
			}
		}

		go func() {
			click := RequestClickCreate(c, ipDB, id, realURL)
			db.Create(click)
			db.Model(&datatypes.Entry{}).Where("id = ?", id).UpdateColumn("count", gorm.Expr("count + 1"))
		}()

		c.Redirect(http.StatusSeeOther, realURL)

	}
}
