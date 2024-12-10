package datatypes

import "time"

type Entry struct {
	ID           int        `gorm:"primaryKey" json:"-"`
	Param        string     `gorm:"index;unique" json:"handle"`
	User         string     `gorm:"index" json:"user"`
	RealURL      string     `json:"url"`
	Custom       bool       `json:"custom"`
	Count        int        `json:"-"`
	Archived     bool       `json:"-"`
	Date         time.Time  `json:"-"`
	ArchivedDate *time.Time `json:"-"`
}
