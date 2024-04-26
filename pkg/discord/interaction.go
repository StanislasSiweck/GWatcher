package discord

import (
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/internal/pkg/session"
	"github.com/bwmarrin/discordgo"
	"log"
)

func UpdateEmbed(guild class.Guild) {
	mes, err := session.DG.ChannelMessage(guild.ChanelID, guild.MessageID)
	if err != nil {
		log.Println("error while fetching message: ", err)
		return
	}

	messageUpdate := guild.Infos.UpdateMessage()
	messageUpdate.ID = mes.ID
	messageUpdate.Channel = mes.ChannelID

	_, err = session.DG.ChannelMessageEditComplex(messageUpdate)
	if err != nil {
		log.Println(err)
	}
}

func FoundGuild(s *discordgo.Session, i *discordgo.InteractionCreate, Guilds map[string]class.Guild) (class.Guild, bool) {
	guild, ok := Guilds[i.GuildID]
	if !ok {
		if err := BasicResponse(s, i, "Guild not found"); err != nil {
			log.Println(err)
		}
		return class.Guild{}, false
	}
	return guild, true
}

func BasicResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
}
