package request

import (
	"bot-serveur-info/internal/pkg/sql"
	"bot-serveur-info/internal/pkg/sql/model"
)

func GetGuildsWithServers() ([]model.Guild, error) {
	var guilds []model.Guild
	err := sql.DB.Preload("Servers").Find(&guilds).Error
	return guilds, err
}

func GetGuildUnScoped(guild model.Guild) model.Guild {
	var g model.Guild
	sql.DB.Unscoped().Where("id = ?", guild.ID).First(&g)
	return g
}

func AddGuild(guild model.Guild) error {
	return sql.DB.Create(&guild).Error
}

func UpdateGuild(guild model.Guild) error {
	return sql.DB.Updates(&guild).Error
}

func RemoveGuild(guidId string) error {
	return sql.DB.Where("id = ?", guidId).Delete(&model.Guild{}).Error
}

func RestoredGuild(guild model.Guild) error {
	return sql.DB.Unscoped().Model(&model.Guild{}).Where("id = ?", guild.ID).Update("deleted_at", nil).Error
}
