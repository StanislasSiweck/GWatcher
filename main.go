package main

import (
	"bot-serveur-info/internal/discord"
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/internal/pkg/session"
	"bot-serveur-info/internal/pkg/sql"
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
	"bot-serveur-info/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
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
			newGuild := class.InitGuild(guild.ID, guild.ChannelID, guild.MessageID, false)
			newGuild.SetDisplayInfo(class.NewDisplay(guild.Servers, 0))
			discord.Guilds[strconv.Itoa(int(guild.ID))] = newGuild
		}
	}

	err := appCommands()
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
				guild := class.InitGuild(0, mes.ChannelID, mes.ID, true)
				guild.SetDisplayInfo(class.NewDisplay([]model.Server{}, 0))
				discord.Guilds[os.Getenv("DISCORD_GUILD_ID")] = guild
				return
			}
		}

		logger.Fatal("The channel ID is not in the guild")
	}
}

func appCommands() error {
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

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "message",
					Description: "Create basic message",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
				{
					Name:        "add",
					Description: "Add a server",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "ip",
							Description: "Server IP",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "port",
							Description: "Server port",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
				{
					Name:        "remove",
					Description: "Remove a server",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "ip",
							Description: "Server IP",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "port",
							Description: "Server port",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
					},
				},
			},
			Description: "Add or remove a server",
		},
	}
)
