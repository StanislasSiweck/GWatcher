package sql

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
}

// AddServer Add a new server to the database
func AddServer(server Server) error {
	return DB.Create(&server).Error
}

// RemoveServer Remove a server from the database by ip and port
func RemoveServer(ip, port string) error {
	return DB.Where("ip = ? AND port = ?", ip, port).Delete(&Server{}).Error
}
