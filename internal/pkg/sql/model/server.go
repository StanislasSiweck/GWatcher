package model

import (
	"gorm.io/gorm"
	"time"
)

type Server struct {
	ID uint `gorm:"primaryKey"`

	IP   string
	Port string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	GuildID *uint
}
