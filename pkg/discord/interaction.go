package discord

import (
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/internal/pkg/session"
	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"
	"log"
	"log/slog"
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
		slog.Error("Can't edit message", tint.Err(err), "message_id", mes.ID, "channel_id", mes.ChannelID, "guild_id")
	}
}

func FoundGuild(s *discordgo.Session, i *discordgo.InteractionCreate, Guilds map[string]class.Guild) (class.Guild, bool) {
	guild, ok := Guilds[i.GuildID]
	if !ok {
		if err := BasicResponse(s, i, "Guild not found"); err != nil {
			slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
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
