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

	GuildID *uint
}

func AddServer(server Server) error {
	return DB.Create(&server).Error
}

func GetServerUnScoped(server Server) Server {
	var s Server
	DB.Unscoped().Where("ip = ? AND port = ? AND guild_id = ?", server.IP, server.Port, server.GuildID).First(&s)
	return s
}

func RestoredServer(server Server) error {
	return DB.Unscoped().Model(&Server{}).Where("ip = ? AND port = ? AND guild_id = ?", server.IP, server.Port, server.GuildID).Update("deleted_at", nil).Error
}

func RemoveServer(ip, port string, guidId uint) error {
	return DB.Where("ip = ? AND port = ? AND guild_id = ?", ip, port, guidId).Delete(&Server{}).Error
}
