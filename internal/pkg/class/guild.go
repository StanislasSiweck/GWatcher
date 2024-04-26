package class

import (
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
)

type Guild struct {
	guildId   uint
	ChanelID  string
	MessageID string

	Infos   DisplayInfo
	isLocal bool
}

func CreateGuild(guildId uint, channelID, messageID string) Guild {
	guild := InitGuild(guildId, channelID, messageID, false)
	if err := guild.CreateGuild(); err != nil {
		return Guild{}
	}
	return guild
}

func InitGuild(guildId uint, channelID, messageID string, isLocal bool) Guild {
	return Guild{
		guildId:   guildId,
		ChanelID:  channelID,
		MessageID: messageID,
		isLocal:   isLocal,
	}
}

func (t *Guild) SetDisplayInfo(dInfo DisplayInfo) {
	t.Infos = dInfo
}

func (t *Guild) NextPage() {
	t.Infos.NextPage()
}

func (t *Guild) PreviousPage() {
	t.Infos.PreviousPage()
}

func (t *Guild) CreateGuild() error {
	if t.isLocal {
		return nil
	}
	guild := model.Guild{
		ID:        t.guildId,
		ChannelID: t.ChanelID,
		MessageID: t.MessageID,
	}

	if request.GetGuildUnScoped(guild).ID != 0 {
		return request.RestoredGuild(guild)
	}

	return request.AddGuild(guild)
}

func (t *Guild) UpdateGuild() error {
	if t.isLocal {
		return nil
	}
	guild := model.Guild{
		ID:        t.guildId,
		ChannelID: t.ChanelID,
		MessageID: t.MessageID,
	}
	return request.UpdateGuild(guild)
}

func (t *Guild) AddServer(server model.Server) (DisplayInfo, error) {
	t.Infos.AddServer(server)
	server.GuildID = &t.guildId

	if t.isLocal {
		return t.Infos, nil
	}

	if request.GetServerUnScoped(server).IP != "" {
		return t.Infos, request.RestoredServer(server)
	}
	return t.Infos, request.AddServer(server)
}

func (t *Guild) RemoveServer(server model.Server) (DisplayInfo, error) {
	t.Infos.RemoveServer(server)
	if t.isLocal {
		return t.Infos, nil
	}
	return t.Infos, request.RemoveServer(server.IP, server.Port, t.guildId)
}

func (t *Guild) HasServer(IP, Port string) bool {
	return t.Infos.HasServer(IP, Port)
}
