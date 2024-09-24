package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"redir/datatypes"
	"strings"
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

func ParseCustomStruct(jsonStr string) (*datatypes.Custom, error) {
	var custom datatypes.Custom
	err := json.Unmarshal([]byte(jsonStr), &custom)
	if err != nil {
		return nil, errors.New("failed to parse JSON")
	}

	return &custom, nil
}

func GetRealURL(db *gorm.DB, id int) (string, error) {
	var realURL string
	result := db.Table("entries").Select("real_url").Where("id = ? AND archived = ?", id, false).Scan(&realURL)
	if result.Error != nil {
		return "", result.Error
	}
	return realURL, nil
}

func GetRealURLAndUserByCustom(db *gorm.DB, custom string) (realURL, user string, err error) {
	var entry datatypes.Entry
	result := db.Where("custom_handle = ? AND archived = ?", custom, false).First(&entry)
	if result.Error != nil {
		err = result.Error
		return
	}
	realURL = entry.RealURL
	user = entry.User
	return
}

func CheckPaymentStatus(userid string, httpClient *http.Client) (bool, error) {
	checkURL := os.Getenv("PAY_API_URL")
	if checkURL == "" {
		checkURL = "https://pay.shortentrack.com"
	}

	url := fmt.Sprintf("%s/check/%s", checkURL, userid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	passcode := os.Getenv("CHECK_PASSCODE")
	if passcode == "" {
		return false, errors.New("missing passcode")
	}
	req.Header.Set("X-Passcode-ID", passcode)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 || resp.StatusCode == 500 {
		return false, errors.New("server error")
	}

	var result struct {
		ID     string `json:"id"`
		Paying bool   `json:"paying"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.Paying, nil
}

func RequestClickCreate(c *gin.Context, ipDB *geoip2.Reader, id int, realURL string) *datatypes.Click {
	var city string
	var country string

	ipStr := c.ClientIP()

	if ipStr == "" || ipStr == "::1" {
		ipStr = c.Request.Header.Get("X-Forwarded-For")
	}

	if ipStr != "" {
		if commaIndex := strings.Index(ipStr, ","); commaIndex != -1 {
			ipStr = ipStr[:commaIndex]
		}

		ip := net.ParseIP(ipStr)
		if ip != nil {
			record, err := ipDB.City(ip)
			if err == nil && record != nil {
				city = record.City.Names["en"]
				country = record.Country.Names["en"]
			}
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
