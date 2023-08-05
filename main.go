package main

import (
	"bot-serveur-info/discord"
	"bot-serveur-info/serveur"
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
			discord.AllServers[s.IP+":"+s.Port] = s
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

	channelID := os.Getenv("DISCORD_CHANEL_ID")
	message := "🤔"

	mes, err := discord.DG.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Fatal("Error sending message :", err)
	}

	go serveur.GetServerInfo(mes)

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutting down")
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "addserver",
			Description: "Add a server to the list",
			Type:        discordgo.ChatApplicationCommand,
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
			Name:        "removeserver",
			Description: "Remove a server from the list",
			Type:        discordgo.ChatApplicationCommand,
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
	}
)
