package datatypes

import (
	"time"
)

type Click struct {
	ID       int `gorm:"primaryKey"`
	ParamKey int `gorm:"index"`
	Time     time.Time
	City     string
	Country  string
	Browser  string
	OS       string
	Platform string
	Mobile   bool
	Bot      bool
	FromQR   bool
}

type Entry struct {
	ID       int `gorm:"primaryKey"`
	User     string
	RealURL  string
	Archived bool
	Date     time.Time
}
