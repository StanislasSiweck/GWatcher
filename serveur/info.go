package serveur

import (
	"bot-serveur-info/discord"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/rumblefrog/go-a2s"
	"log"
	"time"
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

// GetServerInfo is a function that continuously fetches
// server info from all registered Discord servers and sends
// updates to the appropriate Discord channel every minute.
//
// The information gathered includes the server's IP address,
// port, whether it is password protected, the number of
// players currently connected, the maximum number of
// players, the map being played, and a direct connection
// link using the Steam protocol.
//
// It takes a Message object from the discordgo library as
// input, which is used for message editing to send the
// updates to the Discord channel.
//
// Note: This function contains an infinite loop
// and should be run as a goroutine.
//
// Parameters:
// mes (*discordgo.Message): The Discord message object used for message editing.
//
// Returns:
// There is no return value for this function. All errors are logged and handled within the function.
//
// And remember, due to the loop structure in this function,
// any fatal errors encountered will stop the process.
func GetServerInfo(mes *discordgo.Message) {
	for {
		var Fields []*discordgo.MessageEmbedField
		for _, server := range discord.AllServers {

			// initiating a new All-2-Steam client with the server's IP and port
			client, err := a2s.NewClient(server.IP + ":" + server.Port)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			// querying the server info from the client
			info, err := client.QueryInfo()
			if err != nil {
				log.Println(err.Error())
				continue
			}

			// determines whether the server is password-protected
			isPassword := "🔓"
			if info.Visibility {
				isPassword = "🔒"
			}

			// creates the message Discord field's value
			value := fmt.Sprintf("👥 ┃ Players connected `%v/%v` \n", info.Players, info.MaxPlayers)
			value += fmt.Sprintf("🌍 ┃ Map `%v` \n", info.Map)
			value += fmt.Sprintf("📡 ┃ **steam://connect/%v:%v**", server.IP, server.Port)

			// creates a new Discord message field with the specified name and value
			Field := &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("%v %v", info.Name, isPassword),
				Value: value,
			}

			// adds the new field to the Fields slice
			Fields = append(Fields, Field)
		}

		// edit a Discord message with the specified fields
		content := ""
		test := discordgo.MessageEdit{
			Content: &content,
			ID:      mes.ID,
			Channel: mes.ChannelID,
			Embed: &discordgo.MessageEmbed{
				Title:       "Server watch list",
				Description: "━━━━━━━━━━━━━━━━━━━━━━━━",
				Color:       0x5ad65c,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "Update",
				},
				Timestamp: time.Now().Format(time.RFC3339),
				Fields:    Fields,
			},
		}
		_, err := discord.DG.ChannelMessageEditComplex(&test)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Minute)
	}
}
