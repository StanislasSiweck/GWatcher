package main

import (
	"bot-serveur-info/discord"
	"bot-serveur-info/sql"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	discord.NewAuth()
	defer discord.DG.Close()

	if err := sql.ConnectDB(); err == nil {
		if err := sql.Migrate(); err != nil {
			log.Fatal(err)
		}

		guilds, err := sql.GetGuildsWithServers()
		if err != nil {
			log.Fatal(err)
		}

		for _, guild := range guilds {
			newGuild := discord.NewGuild(guild.ID, guild.ChannelID, guild.MessageID, false)
			newGuild.SetDisplayInfo(discord.NewDisplay(guild.Servers, 0))
			discord.Guilds[strconv.Itoa(int(guild.ID))] = newGuild
		}
	}

	err := appCommands()
	if err != nil {
		log.Fatal(err)
	}

	discord.DG.AddHandler(discord.InteractionCreate)

	BDD, _ := sql.DB.DB()
	if BDD.Ping() == nil {
		discord.DG.AddHandler(discord.GuildCreate)
		discord.DG.AddHandler(discord.GuildDelete)
	} else {
		localUsage(err)
	}

	go discord.RefreshServerInfo()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

func localUsage(err error) {
	var mes *discordgo.Message
	messageId := os.Getenv("DISCORD_MESSAGE_ID")

	if messageId != "" {
		mes, err = discord.DG.ChannelMessage(os.Getenv("DISCORD_CHANEL_ID"), messageId)
		if err != nil {
			log.Fatal("Error getting message :", err)
		}

		if mes.Author.ID != discord.DG.State.User.ID {
			mes = nil
		}
	}

	if mes == nil {
		mes, err = discord.DG.ChannelMessageSend(os.Getenv("DISCORD_CHANEL_ID"), "🤔")
		if err != nil {
			log.Fatal("Error sending message :", err)
		}
	}

	if mes != nil {
		guild := discord.NewGuild(0, mes.ChannelID, mes.ID, true)
		guild.SetDisplayInfo(discord.NewDisplay([]sql.Server{}, 0))
		discord.Guilds[os.Getenv("DISCORD_GUILD_ID")] = guild
	}
}

func appCommands() error {
	existingCommands, err := discord.DG.ApplicationCommands(discord.DG.State.User.ID, "")
	if err != nil {
		return err
	}

	for _, command := range existingCommands {
		if err := discord.DG.ApplicationCommandDelete(discord.DG.State.User.ID, "", command.ID); err != nil {
			return err
		}
	}

	for _, command := range commands {
		_, err = discord.DG.ApplicationCommandCreate(discord.DG.State.User.ID, "", command)
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
