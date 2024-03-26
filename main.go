package main

import (
	"bot-serveur-info/discord"
	"bot-serveur-info/sql"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

func main() {
	discord.NewAuth()
	defer discord.DG.Close()

	if err := sql.ConnectDB(); err == nil {
		if err := sql.Migrate(); err != nil {
			log.Fatal(err)
		}

		var server []sql.Server
		sql.DB.Find(&server)

		for _, s := range server {
			discord.AllServers = append(discord.AllServers, s)
		}
	}

	appC, err := discord.DG.ApplicationCommands(discord.DG.State.User.ID, "")

	for _, command := range appC {
		if err := discord.DG.ApplicationCommandDelete(discord.DG.State.User.ID, "", command.ID); err != nil {
			log.Fatal(err)
		}
	}

	for _, command := range commands {
		_, err = discord.DG.ApplicationCommandCreate(discord.DG.State.User.ID, "", command)
		if err != nil {
			log.Fatal(err)
		}
	}

	discord.DG.AddHandler(discord.InteractionCreate)

	messageId := os.Getenv("DISCORD_MESSAGE_ID")

	if messageId != "" {
		discord.Mes, err = discord.DG.ChannelMessage(os.Getenv("DISCORD_CHANEL_ID"), messageId)
		if err != nil {
			log.Fatal("Error getting message :", err)
		}
		if discord.Mes.Author.ID != discord.DG.State.User.ID {
			discord.Mes = nil
		}
	}

	if discord.Mes == nil {
		discord.Mes, err = discord.DG.ChannelMessageSend(os.Getenv("DISCORD_CHANEL_ID"), "🤔")
		if err != nil {
			log.Fatal("Error sending message :", err)
		}
	}

	go discord.RefreshServerInfo()

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "server",
			Options: []*discordgo.ApplicationCommandOption{
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
