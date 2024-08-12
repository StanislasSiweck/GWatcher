package controller

import "bot-serveur-info/internal/pkg/sql/model"

type Guild interface {
	SetDisplayInfo(dInfo DisplayInfo)
	NextPage()
	PreviousPage()
	CreateGuild() error
	UpdateGuild() error
	AddServer(server model.Server) (DisplayInfo, error)
	RemoveServer(server model.Server) (DisplayInfo, error)
	HasServer(IP, Port string) bool
	Infos() DisplayInfo
	ChangeMessage(chanelID, messageID string)
	Message() (string, string)
}
