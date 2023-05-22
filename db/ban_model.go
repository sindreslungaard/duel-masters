package db

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type BanType string

const (
	UserBan BanType = "user"
	IPBan   BanType = "ip"
)

func (b *BanType) Scan(value interface{}) error {
	*b = BanType(value.([]byte))
	return nil
}

func (b BanType) Value() (driver.Value, error) {
	return string(b), nil
}

type Ban struct {
	gorm.Model
	ID      uint
	BanType BanType `sql:"type:ENUM('user', 'ip')"`
	Value   string
	Reason  string
	Created uint
	Expires uint
}
