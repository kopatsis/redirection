package datatypes

import (
	"time"
)

type Click struct {
	ID        int `gorm:"primaryKey"`
	ParamKey  int `gorm:"index"`
	Time      time.Time
	RealURL   string
	City      string
	Country   string
	Browser   string
	OS        string
	Platform  string
	Mobile    bool
	Bot       bool
	FromQR    bool
	IPAddress string
}
