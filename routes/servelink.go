package routes

import (
	"net/http"
	"redir/convert"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/gorm"
)

func Redirect(db *gorm.DB, ipDB *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		id, err := convert.FromSixFour(param)
		if err != nil {
			MainRedirect(c)
			return
		}

		realURL, err := GetRealURL(db, id)
		if err != nil {
			MainRedirect(c)
			return
		}

		go func() {
			click := RequestClickCreate(c, ipDB, id)
			db.Create(click)
		}()

		c.Redirect(http.StatusFound, realURL)

	}
}
