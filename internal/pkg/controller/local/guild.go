package local

import (
	"bot-serveur-info/internal/pkg/controller"
	"bot-serveur-info/internal/pkg/sql/model"
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
	return nil
}

func (t *Guild) UpdateGuild() error {
	return nil
}

func (t *Guild) AddServer(server model.Server) (controller.DisplayInfo, error) {
	t.infos.AddServer(server)
	server.GuildID = &t.guildId

	return t.infos, nil
}

func (t *Guild) RemoveServer(server model.Server) (controller.DisplayInfo, error) {
	t.infos.RemoveServer(server)
	return t.infos, nil
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
