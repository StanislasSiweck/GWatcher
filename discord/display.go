package discord

import (
	"bot-serveur-info/serveur"
	"bot-serveur-info/sql"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

type DisplayInfo struct {
	servers []sql.Server
	page    int
}

type Server struct {
	//Name *string
	IP   string
	Port string
}

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

func NewDisplay(servers []sql.Server, page int) DisplayInfo {
	return DisplayInfo{
		servers: servers,
		page:    page,
	}
}

func (d *DisplayInfo) maxPage() int {
	return (len(d.servers) / 2) + len(d.servers)%2
}

func (d *DisplayInfo) AddServer(server sql.Server) {
	d.servers = append(d.servers, sql.Server{
		IP:   server.IP,
		Port: server.Port,
	})
}

func (d *DisplayInfo) RemoveServer(server sql.Server) {
	for i, s := range d.servers {
		if s.IP == server.IP && s.Port == server.Port {
			d.servers = append(d.servers[:i], d.servers[i+1:]...)
			break
		}
	}
}

func (d *DisplayInfo) HasServer(IP, Port string) bool {
	for _, s := range d.servers {
		if s.IP == IP && s.Port == Port {
			return true
		}
	}
	return false
}

func (d *DisplayInfo) Fields() (Fields []*discordgo.MessageEmbedField) {
	for count, server := range d.servers {
		if count < d.page*2 || count > d.page*2+1 {
			continue
		}

		var field *discordgo.MessageEmbedField
		info, err := serveur.GetServerInfo(server)
		if err != nil {
			field = &discordgo.MessageEmbedField{
				Name:  "Error",
				Value: "Error while fetching server info",
			}
		} else {
			field = serveur.CreateField(info, server)
		}

		Fields = append(Fields, field)
	}

	return Fields
}

func (d *DisplayInfo) UpdateMessage() *discordgo.MessageEdit {
	fields := d.Fields()

	left := constLeft
	if d.page == 0 {
		left.Disabled = true
	}

	right := constRight
	if d.page == d.maxPage()-1 {
		right.Disabled = true
	}

	content := ""
	messageEdit := discordgo.MessageEdit{
		Content: &content,
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
			Title:       "Server watch list (Page " + strconv.Itoa(d.page+1) + "/" + strconv.Itoa(d.maxPage()) + ")",
			Description: "━━━━━━━━━━━━━━━━━━━━━━━━",
			Color:       0x5ad65c,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "UpdateGuild",
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Fields:    fields,
		},
	}

	return &messageEdit
}
