package discord

import (
	"bot-serveur-info/serveur"
	"bot-serveur-info/sql"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

var AllServers = []sql.Server{}
var Mes *discordgo.Message
var page = 0

var (
	constRight = discordgo.Button{
		Emoji: discordgo.ComponentEmoji{
			Name: "➡️",
		},
		Style:    discordgo.PrimaryButton,
		CustomID: "right",
	}
	constLeft = discordgo.Button{
		Emoji: discordgo.ComponentEmoji{
			Name: "⬅️",
		},
		Style:    discordgo.PrimaryButton,
		CustomID: "left",
	}
)

var (
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"server add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			addServerCommand(s, i)
		},
		"server remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			removeServerCommand(s, i)
		},
	}
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Refreshed",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			ServerInfo()
		},
		"right": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			page++
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Page " + strconv.Itoa(page+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			ServerInfo()
		},
		"left": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if page == 0 {
				return
			}
			page--
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Page " + strconv.Itoa(page+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			ServerInfo()
		},
	}
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		command := data.Name

		if len(data.Options) == 1 {
			command = data.Name + " " + data.Options[0].Name
		}
		if h, ok := commandsHandlers[command]; ok {
			h(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
	}
}

func addServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	appOption := i.ApplicationCommandData().Options[0]
	server := sql.Server{
		IP:   appOption.Options[0].StringValue(),
		Port: appOption.Options[1].StringValue(),
	}

	if err := sql.AddServer(server); err != nil { // Create the server in the database
		log.Println(err)
	}

	AllServers = append(AllServers, server) // Add to local list

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server added",
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func removeServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	appOption := i.ApplicationCommandData().Options[0]

	for i, server := range AllServers {
		if server.IP == appOption.Options[0].StringValue() && server.Port == appOption.Options[1].StringValue() {
			AllServers = append(AllServers[:i], AllServers[i+1:]...)
			break
		}
	}
	ip, port := appOption.Options[0].StringValue(), appOption.Options[1].StringValue()

	if err := sql.RemoveServer(ip, port); err != nil { // Remove from database
		log.Println(err)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server removed",
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func RefreshServerInfo() {
	for {
		ServerInfo()
		time.Sleep(1 * time.Minute)
	}
}

func ServerInfo() {
	var Fields []*discordgo.MessageEmbedField
	for count, server := range AllServers {
		if count < page*2 || count > page*2+1 {
			continue
		}

		var field *discordgo.MessageEmbedField
		info, err := serveur.GetServerInfo(server)
		if err != nil {
			log.Println(err)
			field = &discordgo.MessageEmbedField{
				Name:  "Error",
				Value: "Error while fetching server info",
			}
		} else {
			field = serveur.CreateField(info, server)
		}

		Fields = append(Fields, field)

	}

	maxPage := len(AllServers) / 2

	// check modulo
	if len(AllServers)%2 == 1 {
		maxPage++
	}

	left := constLeft
	if page == 0 {
		left.Disabled = true
	}

	right := constRight
	if page == maxPage-1 {
		right.Disabled = true
	}

	// edit a Discord message with the specified fields
	content := ""
	messageEdit := discordgo.MessageEdit{
		Content: &content,
		ID:      Mes.ID,
		Channel: Mes.ChannelID,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					left,
					discordgo.Button{
						Emoji: discordgo.ComponentEmoji{
							Name: "🔄",
						},
						Style:    discordgo.PrimaryButton,
						CustomID: "update",
					},
					right,
				},
			},
		},
		Embed: &discordgo.MessageEmbed{
			Title:       "Server watch list (Page " + strconv.Itoa(page+1) + "/" + strconv.Itoa(maxPage) + ")",
			Description: "━━━━━━━━━━━━━━━━━━━━━━━━",
			Color:       0x5ad65c,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Update",
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Fields:    Fields,
		},
	}
	_, err := DG.ChannelMessageEditComplex(&messageEdit)
	if err != nil {
		log.Fatal(err)
	}
}
