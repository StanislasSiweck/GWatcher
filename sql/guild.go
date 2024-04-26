package sql

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

func AddGuild(guild Guild) error {
	return DB.Create(&guild).Error
}

func UpdateGuild(guild Guild) error {
	return DB.Updates(&guild).Error
}

func GetGuildUnScoped(guild Guild) Guild {
	var g Guild
	DB.Unscoped().Where("id = ?", guild.ID).First(&g)
	return g
}

func RestoredGuild(guild Guild) error {
	return DB.Unscoped().Model(&Guild{}).Where("id = ?", guild.ID).Update("deleted_at", nil).Error
}

func RemoveGuild(guidId string) error {
	return DB.Where("id = ?", guidId).Delete(&Guild{}).Error
}

func GetGuildsWithServers() ([]Guild, error) {
	var guilds []Guild
	err := DB.Preload("Servers").Find(&guilds).Error
	return guilds, err
}
