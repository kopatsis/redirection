package datatypes

import "time"

type Entry struct {
	ID           int `gorm:"primaryKey"`
	User         string
	RealURL      string
	Archived     bool
	Date         time.Time
	ArchivedDate time.Time
}
