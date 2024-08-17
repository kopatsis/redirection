package datatypes

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Click struct {
	gorm.Model
	ParamKey int
	Time     pq.NullTime
	City     string
	Country  string
	Browser  string
	OS       string
	Platform string
	Mobile   bool
	Bot      bool
	FromQR   bool
}
