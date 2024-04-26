package class

import (
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/pkg/serveur"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

type DisplayInfo struct {
	servers []model.Server
	Page    int
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

func NewDisplay(servers []model.Server, page int) DisplayInfo {
	return DisplayInfo{
		servers: servers,
		Page:    page,
	}
}

func (d *DisplayInfo) NextPage() {
	if d.Page < d.MaxPage()-1 {
		d.Page++
	}
}

func (d *DisplayInfo) PreviousPage() {
	if d.Page > 0 {
		d.Page--
	}
}

func (d *DisplayInfo) MaxPage() int {
	return (len(d.servers) / 2) + len(d.servers)%2
}

func (d *DisplayInfo) AddServer(server model.Server) {
	d.servers = append(d.servers, model.Server{
		IP:   server.IP,
		Port: server.Port,
	})
}

func (d *DisplayInfo) RemoveServer(server model.Server) {
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
		if count < d.Page*2 || count > d.Page*2+1 {
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
	if d.Page == 0 {
		left.Disabled = true
	}

	right := constRight
	if d.Page == d.MaxPage()-1 {
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
			Title:       "Server watch list (Page " + strconv.Itoa(d.Page+1) + "/" + strconv.Itoa(d.MaxPage()) + ")",
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
