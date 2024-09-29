package routes

import (
	"context"
	"errors"
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

		var userID string

		retStr, err := rdb.Get(context.Background(), param).Result()
		if (len(retStr) >= 3 && retStr[:3] == ":e:") || retStr == ":a:" {
			MainRedirect(c, true)
			return
		} else if err == nil {
			if custom {
				customStruct, parseErr := ParseCustomStruct(retStr)
				if parseErr != nil {
					err = parseErr
				} else if customStruct == nil || customStruct.URL == "" || customStruct.UserID == "" {
					err = errors.New("improperly formatted custom struct for custom handle")
				} else {
					realURL = customStruct.URL
					userID = customStruct.UserID
				}
			} else {
				realURL = retStr
			}
		}

		if err != nil {
			if custom {
				realURL, userID, err = GetRealURLAndUserByCustom(db, param)
			} else {
				realURL, err = GetRealURL(db, id)
			}

			if err != nil {
				MainRedirect(c, true)
				return
			}
		}

		if custom {
			check, err := CheckUserPaying(rdb, userID)
			if !check || err != nil {
				MainRedirect(c, true)
				return
			}
		}

		if c.GetHeader("Already-Redirected") != "true" {
			go func() {
				click := RequestClickCreate(c, ipDB, id, realURL, param, custom)
				db.Create(click)
				db.Model(&datatypes.Entry{}).Where("id = ?", id).UpdateColumn("count", gorm.Expr("count + 1"))
			}()
		}

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Header("Already-Redirected", "true")

		c.Redirect(http.StatusSeeOther, realURL)

	}
}
