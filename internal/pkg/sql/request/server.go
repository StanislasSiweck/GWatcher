package request

import (
	"bot-serveur-info/internal/pkg/sql"
	"bot-serveur-info/internal/pkg/sql/model"
)

func AddServer(server model.Server) error {
	return sql.DB.Create(&server).Error
}

func RemoveServer(ip, port string, guidId uint) error {
	return sql.DB.Where("ip = ? AND port = ? AND guild_id = ?", ip, port, guidId).Delete(&model.Server{}).Error
}

func GetServerUnScoped(server model.Server) model.Server {
	var s model.Server
	sql.DB.Unscoped().Where("ip = ? AND port = ? AND guild_id = ?", server.IP, server.Port, server.GuildID).First(&s)
	return s
}

func RestoredServer(server model.Server) error {
	return sql.DB.Unscoped().Model(&model.Server{}).Where("ip = ? AND port = ? AND guild_id = ?", server.IP, server.Port, server.GuildID).Update("deleted_at", nil).Error
}
