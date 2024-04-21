package discord

import (
	"bot-serveur-info/sql"
)

type Guild struct {
	guildId   uint
	chanelID  string
	messageID string

	infos   DisplayInfo
	isLocal bool
}

func CreateGuild(guildId uint, channelID, messageID string) Guild {
	guild := NewGuild(guildId, channelID, messageID, false)
	err := guild.CreateGuild()
	if err != nil {
		return Guild{}
	}
	return guild
}

func NewGuild(guildId uint, channelID, messageID string, isLocal bool) Guild {
	return Guild{
		guildId:   guildId,
		chanelID:  channelID,
		messageID: messageID,
		isLocal:   isLocal,
	}
}

func (t *Guild) SetDisplayInfo(dInfo DisplayInfo) {
	t.infos = dInfo
}

func (t *Guild) CreateGuild() error {
	if t.isLocal {
		return nil
	}
	guild := sql.Guild{
		ID:        t.guildId,
		ChannelID: t.chanelID,
		MessageID: t.messageID,
	}

	if sql.GetGuildUnScoped(guild).ID != 0 {
		return sql.RestoredGuild(guild)
	}

	return sql.AddGuild(guild)
}

func (t *Guild) UpdateGuild() error {
	if t.isLocal {
		return nil
	}
	guild := sql.Guild{
		ID:        t.guildId,
		ChannelID: t.chanelID,
		MessageID: t.messageID,
	}
	return sql.UpdateGuild(guild)
}

func (t *Guild) AddServer(server sql.Server) (DisplayInfo, error) {
	t.infos.AddServer(server)
	server.GuildID = &t.guildId

	if t.isLocal {
		return t.infos, nil
	}

	if sql.GetServerUnScoped(server).IP != "" {
		return t.infos, sql.RestoredServer(server)
	}
	return t.infos, sql.AddServer(server)
}

func (t *Guild) RemoveServer(server sql.Server) (DisplayInfo, error) {
	t.infos.RemoveServer(server)
	if t.isLocal {
		return t.infos, nil
	}
	return t.infos, sql.RemoveServer(server.IP, server.Port, t.guildId)
}

func (t *Guild) HasServer(IP, Port string) bool {
	return t.infos.HasServer(IP, Port)
}
