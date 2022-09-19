package models

import "time"

type Address struct {
	ID         string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	User       User
	UserID     string `gorm:"size:36;index"`
	Name       string `gorm:"size100"`
	IsPrimary  bool
	CityID     string `gorm:"size100"`
	ProvinceID string `gorm:"size100"`
	Address1   string `gorm:"size255"`
	Address2   string `gorm:"size255"`
	Phone      string `gorm:"size100"`
	Email      string `gorm:"size100"`
	PostCode   string `gorm:"size100"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
