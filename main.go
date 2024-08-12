package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"

	"bot-serveur-info/internal/discord"
	"bot-serveur-info/internal/pkg/controller"
	csql "bot-serveur-info/internal/pkg/controller/bdd"
	clocal "bot-serveur-info/internal/pkg/controller/local"
	"bot-serveur-info/internal/pkg/session"
	"bot-serveur-info/internal/pkg/sql"
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
	pdiscord "bot-serveur-info/pkg/discord"
	"bot-serveur-info/pkg/logger"
)

func main() {
	logger.New()

	session.NewAuth()
	defer session.DG.Close()

	if err := sql.ConnectDB(); err == nil {
		if err := sql.Migrate(); err != nil {
			logger.Fatal("Migration goes wrong", tint.Err(err))
		}

		guilds, err := request.GetGuildsWithServers()
		if err != nil {
			logger.Fatal("Guilds cannot recover", tint.Err(err))
		}

		for _, guild := range guilds {
			newGuild := csql.InitGuild(guild.ID, guild.ChannelID, guild.MessageID)
			newGuild.SetDisplayInfo(controller.NewDisplay(guild.Servers, 0))
			discord.Guilds[strconv.Itoa(int(guild.ID))] = &newGuild
		}
	}

	err := pdiscord.AppCommands()
	if err != nil {
		logger.Fatal("The commands could not be modified", tint.Err(err))
	}

	session.DG.AddHandler(discord.InteractionCreate)

	BDD, _ := sql.DB.DB()
	if BDD.Ping() == nil {
		session.DG.AddHandler(discord.GuildCreate)
		session.DG.AddHandler(discord.GuildDelete)
	} else {
		localUsage(err)
	}

	go discord.RefreshServerInfo()

	slog.Info("Bot is now running.  Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	slog.Info("Gracefully shutting down")
}

func localUsage(err error) {
	if os.Getenv("STORAGE_TYPE") != "local" {
		logger.Fatal("The storage is not local")
	}
	guildId := os.Getenv("DISCORD_GUILD_ID")
	if guildId == "" {
		logger.Fatal("The guild ID is not set")
	}

	var mes *discordgo.Message
	messageId := os.Getenv("DISCORD_MESSAGE_ID")

	if messageId != "" {
		mes, err = session.DG.ChannelMessage(os.Getenv("DISCORD_CHANEL_ID"), messageId)
		if err != nil {
			logger.Fatal("Can't get message in local", tint.Err(err))
		}

		if mes.Author.ID != session.DG.State.User.ID {
			mes = nil
		}
	}

	if mes == nil {
		mes, err = session.DG.ChannelMessageSend(os.Getenv("DISCORD_CHANEL_ID"), "🤔")
		if err != nil {
			logger.Fatal("Can't send message in local", tint.Err(err))
		}
	}

	if mes != nil {
		channels, err := session.DG.GuildChannels(guildId)
		if err != nil {
			logger.Fatal("Can't get guild in local", tint.Err(err))
		}

		for _, channel := range channels {
			if channel.ID == os.Getenv("DISCORD_CHANEL_ID") {
				guild := clocal.InitGuild(0, mes.ChannelID, mes.ID)
				guild.SetDisplayInfo(controller.NewDisplay([]model.Server{}, 0))
				discord.Guilds[os.Getenv("DISCORD_GUILD_ID")] = &guild
				return
			}
		}

		logger.Fatal("The channel ID is not in the guild")
	}
}
