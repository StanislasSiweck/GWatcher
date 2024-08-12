package bdd

import (
	"bot-serveur-info/internal/pkg/controller"
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
)

type Guild struct {
	guildId   uint
	chanelID  string
	messageID string

	infos controller.DisplayInfo
}

func CreateGuild(guildId uint, channelID, messageID string) (*Guild, error) {
	guild := InitGuild(guildId, channelID, messageID)
	if err := guild.CreateGuild(); err != nil {
		return nil, err
	}
	return &guild, nil
}

func InitGuild(guildId uint, channelID, messageID string) Guild {
	return Guild{
		guildId:   guildId,
		chanelID:  channelID,
		messageID: messageID,
	}
}

func (t *Guild) SetDisplayInfo(dInfo controller.DisplayInfo) {
	t.infos = dInfo
}

func (t *Guild) NextPage() {
	t.infos.NextPage()
}

func (t *Guild) PreviousPage() {
	t.infos.PreviousPage()
}

func (t *Guild) CreateGuild() error {
	guild := model.Guild{
		ID:        t.guildId,
		ChannelID: t.chanelID,
		MessageID: t.messageID,
	}

	if request.GetGuildUnScoped(guild).ID != 0 {
		return request.RestoredGuild(guild)
	}

	return request.AddGuild(guild)
}

func (t *Guild) UpdateGuild() error {
	guild := model.Guild{
		ID:        t.guildId,
		ChannelID: t.chanelID,
		MessageID: t.messageID,
	}
	return request.UpdateGuild(guild)
}

func (t *Guild) AddServer(server model.Server) (controller.DisplayInfo, error) {
	t.infos.AddServer(server)
	server.GuildID = &t.guildId

	if request.GetServerUnScoped(server).IP != "" {
		return t.infos, request.RestoredServer(server)
	}
	return t.infos, request.AddServer(server)
}

func (t *Guild) RemoveServer(server model.Server) (controller.DisplayInfo, error) {
	t.infos.RemoveServer(server)
	return t.infos, request.RemoveServer(server.IP, server.Port, t.guildId)
}

func (t *Guild) HasServer(IP, Port string) bool {
	return t.infos.HasServer(IP, Port)
}

func (t *Guild) Infos() controller.DisplayInfo {
	return t.infos
}

func (t *Guild) ChangeMessage(chanelID, messageID string) {
	t.chanelID = chanelID
	t.messageID = messageID
}

func (t *Guild) Message() (string, string) {
	return t.chanelID, t.messageID
}
