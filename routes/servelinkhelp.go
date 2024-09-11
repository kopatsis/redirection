package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"redir/datatypes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"
	"gorm.io/gorm"
)

func MainRedirect(c *gin.Context, isErr bool) {
	mainURL := os.Getenv("mainURL")
	if mainURL == "" {
		mainURL = "https://shortentrack.com"
	}
	if isErr {
		mainURL += "?dne=t"
	}
	c.Redirect(http.StatusFound, mainURL)
}

func GetRealURL(db *gorm.DB, id int) (string, error) {
	var realURL string
	result := db.Table("entries").Select("real_url").Where("id = ? AND archived = ?", id, false).Scan(&realURL)
	if result.Error != nil {
		return "", result.Error
	}
	return realURL, nil
}

func GetRealURLByCustom(db *gorm.DB, custom string) (string, error) {
	var realURL string
	result := db.Table("entries").Select("real_url").Where("custom_handle = ? AND archived = ?", custom, false).Scan(&realURL)
	if result.Error != nil {
		return "", result.Error
	}
	return realURL, nil
}

func RequestClickCreate(c *gin.Context, ipDB *geoip2.Reader, id int, realURL string) *datatypes.Click {
	var city string
	var country string

	ipStr := c.ClientIP()
	ip := net.ParseIP(ipStr)
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
		ParamKey:  id,
		Time:      time.Now(),
		RealURL:   realURL,
		City:      city,
		Country:   country,
		Browser:   browser,
		OS:        os,
		Platform:  platform,
		Mobile:    isMobile,
		Bot:       isBot,
		FromQR:    c.Query("q") == "t",
		IPAddress: hex.EncodeToString(sha256.New().Sum([]byte(ipStr))),
	}

	return &click
}
