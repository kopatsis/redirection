package routes

import (
	"net"
	"redir/datatypes"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
)

func Redirect(db *gorm.DB, ipDB *geoip2.Reader) gin.HandlerFunc {
	return func(c *gin.Context) {

		var city string
		var country string

		ip := net.ParseIP(c.ClientIP())
		if ip != nil {
			record, err := ipDB.City(ip)

			if err == nil && record != nil {
				city = record.City.Names["en"]
				country = record.Country.Names["en"]
			}
		}

		parser := uaparser.NewFromSaved()

		ua := c.Request.UserAgent()
		client := parser.Parse(ua)

		browser := client.UserAgent.Family
		// browserVersion := client.UserAgent.ToVersionString()
		os := client.Os.Family
		platform := client.Device.Family

		uaM := user_agent.New(ua)
		isMobile := uaM.Mobile()
		isBot := uaM.Bot()

		click := datatypes.Click{
			City:     city,
			Country:  country,
			Browser:  browser,
			OS:       os,
			Platform: platform,
			Mobile:   isMobile,
			Bot:      isBot,
		}

		c.JSON(200, &click)
	}

}
