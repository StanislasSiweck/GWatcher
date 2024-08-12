package serveur

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/go-a2s"

	"bot-serveur-info/internal/pkg/sql/model"
)

type Player struct {
	Name string `json:"name"`
	Raw  struct {
		Score int     `json:"score"`
		Time  float64 `json:"time"`
	} `json:"raw"`
}

type ServerInfo struct {
	Name       string   `json:"name"`
	Map        string   `json:"map"`
	Password   bool     `json:"password"`
	Gamemode   string   `json:"gamemode"`
	Maxplayers int      `json:"maxplayers"`
	Players    []Player `json:"players"`
	Bots       []Player `json:"bots"`
	Connect    string   `json:"connect"`
	Ping       int      `json:"ping"`
}

func GetServerInfo(server model.Server) (info *a2s.ServerInfo, err error) {
	// initiating a new All-2-Steam client with the server's IP and port
	client, err := a2s.NewClient(server.IP + ":" + server.Port)
	if err != nil {
		return
	}

	defer client.Close()

	// querying the server info from the client
	info, err = client.QueryInfo()
	if err != nil {
		return
	}
	return
}

func CreateField(info *a2s.ServerInfo, server model.Server) *discordgo.MessageEmbedField {
	// determines whether the server is password-protected
	isPassword := "🔓"
	if info.Visibility {
		isPassword = "🔒"
	}

	value := fmt.Sprintf("👥 ┃ Players connected `%v/%v` \n", info.Players, info.MaxPlayers)
	value += fmt.Sprintf("🌍 ┃ Map `%v` \n", info.Map)
	value += fmt.Sprintf("📡 ┃ **steam://connect/%v:%v**", server.IP, server.Port)

	Field := &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("%v %v", info.Name, isPassword),
		Value: value,
	}
	return Field
}
