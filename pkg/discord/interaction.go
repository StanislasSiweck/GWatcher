package discord

import (
	"log"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"

	"bot-serveur-info/internal/pkg/controller"
	"bot-serveur-info/internal/pkg/session"
)

func AppCommands() error {
	existingCommands, err := session.DG.ApplicationCommands(session.DG.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, command := range existingCommands {
		if err := session.DG.ApplicationCommandDelete(session.DG.State.User.ID, "", command.ID); err != nil {
			return err
		}
	}

	for _, command := range commands {
		_, err = session.DG.ApplicationCommandCreate(session.DG.State.User.ID, "", command)
		if err != nil {
			return err
		}
	}
	return err
}

func UpdateEmbed(guild controller.Guild) {
	chanelId, messageId := guild.Message()
	mes, err := session.DG.ChannelMessage(chanelId, messageId)
	if err != nil {
		log.Println("error while fetching message: ", err)
		return
	}

	infos := guild.Infos()
	messageUpdate := infos.UpdateMessage()
	messageUpdate.ID = mes.ID
	messageUpdate.Channel = mes.ChannelID

	_, err = session.DG.ChannelMessageEditComplex(messageUpdate)
	if err != nil {
		slog.Error("Can't edit message", tint.Err(err), "message_id", mes.ID, "channel_id", mes.ChannelID, "guild_id")
	}
}

func FoundGuild(s *discordgo.Session, i *discordgo.InteractionCreate, Guilds map[string]controller.Guild) (controller.Guild, bool) {
	guild, ok := Guilds[i.GuildID]
	if !ok {
		if err := BasicResponse(s, i, "Guild not found"); err != nil {
			slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
		}
		return nil, false
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
