package discord

import (
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/internal/pkg/session"
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
	"bot-serveur-info/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

var Guilds = map[string]class.Guild{}

func GuildCreate(_ *discordgo.Session, g *discordgo.GuildCreate) {
	_, ok := Guilds[g.ID]
	if ok {
		return
	}
	guildID, err := strconv.Atoi(g.ID)
	if err != nil {
		log.Println(err)
		return
	}
	guild := class.CreateGuild(uint(guildID), "", "")
	guild.SetDisplayInfo(class.NewDisplay([]model.Server{}, 0))
	Guilds[g.ID] = guild
}

func GuildDelete(_ *discordgo.Session, g *discordgo.GuildDelete) {
	_, ok := Guilds[g.ID]
	if !ok {
		return
	}

	if err := request.RemoveGuild(g.ID); err != nil {
		log.Println(err)
		return
	}

	delete(Guilds, g.ID)
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		data := i.ApplicationCommandData()
		command := data.Name

		if len(data.Options) == 1 {
			command = data.Name + " " + data.Options[0].Name
		}

		if h, ok := commandsHandlers[command]; ok {
			guild, done := discord.FoundGuild(s, i, Guilds)
			if !done {
				return
			}

			guild = h(s, i, guild)
			if guild.MessageID != "" {
				discord.UpdateEmbed(guild)
				Guilds[i.GuildID] = guild
			}
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
			guild, done := discord.FoundGuild(s, i, Guilds)
			if !done {
				return
			}
			guild = h(s, i, guild)

			discord.UpdateEmbed(guild)
			Guilds[i.GuildID] = guild
		}
	}
}

func setMessage(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
	mes, err := session.DG.ChannelMessageSend(i.ChannelID, "🤔")
	if err != nil {
		log.Println("Error sending message :", err)
		return class.Guild{}
	}

	guild.ChanelID = mes.ChannelID
	guild.MessageID = mes.ID
	if err := guild.UpdateGuild(); err != nil {
		log.Println(err)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Message send"); err != nil {
		log.Println(err)
		return guild
	}

	return guild
}

func addServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
	appOption := i.ApplicationCommandData().Options[0]

	IP, Port := appOption.Options[0].StringValue(), appOption.Options[1].StringValue()

	if guild.HasServer(IP, Port) {
		if err := discord.BasicResponse(s, i, "Server already existed"); err != nil {
			log.Println(err)
		}
		return class.Guild{}
	}

	server := model.Server{
		IP:   IP,
		Port: Port,
	}

	if _, err := guild.AddServer(server); err != nil { // Create the server in the database
		log.Println(err)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Server added"); err != nil {
		log.Println(err)
		return guild
	}

	return guild
}

func removeServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
	appOptions := i.ApplicationCommandData().Options[0]
	ip, port := appOptions.Options[0].StringValue(), appOptions.Options[1].StringValue()

	if !guild.HasServer(ip, port) {
		if err := discord.BasicResponse(s, i, "Server don't exist"); err != nil {
			log.Println(err)
		}
		return class.Guild{}
	}

	if _, err := guild.RemoveServer(model.Server{IP: ip, Port: port}); err != nil { // Remove from database
		log.Println(err)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Server removed"); err != nil {
		log.Println(err)
		return guild
	}

	return guild
}

// RefreshServerInfo refresh the server info every minute for each guild
func RefreshServerInfo() {
	for {
		RefreshServerEmbed()
		time.Sleep(1 * time.Minute)
	}
}

func RefreshServerEmbed() {
	for _, guild := range Guilds {
		discord.UpdateEmbed(guild)
	}
}
