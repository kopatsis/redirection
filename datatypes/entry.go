package datatypes

import "time"

type Entry struct {
	ID           int       `gorm:"primaryKey" json:"-"`
	User         string    `gorm:"index" json:"user"`
	RealURL      string    `json:"url"`
	CustomHandle string    `gorm:"unique;index" json:"-"`
	Count        int       `json:"-"`
	Archived     bool      `json:"-"`
	Date         time.Time `json:"-"`
	ArchivedDate time.Time `json:"-"`
}
