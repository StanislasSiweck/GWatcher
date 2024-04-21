package model

import (
	"gorm.io/gorm"
	"time"
)

type Guild struct {
	ID uint `gorm:"primaryKey;autoIncrement:false"`

	ChannelID string
	MessageID string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Servers []Server `gorm:"constraint:OnDelete:CASCADE"`
}
