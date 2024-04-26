package discord

import (
	"bot-serveur-info/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

var Guilds = map[string]Guild{}

var (
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild{
		"server add": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			return addServerCommand(s, i)
		},
		"server remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			return removeServerCommand(s, i)
		},
		"server message": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			return sendMessage(s, i)
		},
	}
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild{
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			guild, done := foundGuild(s, i)
			if done {
				return Guild{}
			}
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Refreshed",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			return guild
		},
		"right": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			guild, done := foundGuild(s, i)
			if done {
				return Guild{}
			}
			guild.infos.page++
			Guilds[i.GuildID] = guild

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Page " + strconv.Itoa(guild.infos.page+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			return guild
		},
		"left": func(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
			guild, done := foundGuild(s, i)
			if done {
				return Guild{}
			}
			guild.infos.page--
			Guilds[i.GuildID] = guild

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Page " + strconv.Itoa(guild.infos.page+1),
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			}
			_ = s.InteractionRespond(i.Interaction, response)
			return guild
		},
	}
)

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
	guild := CreateGuild(uint(guildID), "", "")
	guild.SetDisplayInfo(NewDisplay([]sql.Server{}, 0))
	Guilds[g.ID] = guild
}

func GuildDelete(_ *discordgo.Session, g *discordgo.GuildDelete) {
	_, ok := Guilds[g.ID]
	if !ok {
		return
	}

	if err := sql.RemoveGuild(g.ID); err != nil {
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
			guild := h(s, i)
			if guild.messageID != "" {
				DisplayServerInfo(guild)
			}
		}
	case discordgo.InteractionMessageComponent:
		if h, ok := componentsHandlers[i.MessageComponentData().CustomID]; ok {
			guild := h(s, i)
			if guild.messageID != "" {
				DisplayServerInfo(guild)
			}
		}
	}
}

func sendMessage(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
	mes, err := DG.ChannelMessageSend(i.ChannelID, "🤔")
	if err != nil {
		log.Fatal("Error sending message :", err)
	}

	guild, done := foundGuild(s, i)
	if done {
		log.Println("Guild not found")
		return Guild{}
	}
	guild.chanelID = mes.ChannelID
	guild.messageID = mes.ID
	if err := guild.UpdateGuild(); err != nil {
		log.Println(err)
		return Guild{}
	}
	Guilds[i.GuildID] = guild
	return guild
}

func addServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
	guild, done := foundGuild(s, i)
	if done {
		return Guild{}
	}

	appOption := i.ApplicationCommandData().Options[0]

	IP, Port := appOption.Options[0].StringValue(), appOption.Options[1].StringValue()

	if guild.HasServer(IP, Port) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Server already existed",
			},
		})
		if err != nil {
			log.Println(err)
		}
		return Guild{}
	}

	server := sql.Server{
		IP:   IP,
		Port: Port,
	}

	if _, err := guild.AddServer(server); err != nil { // Create the server in the database
		log.Println(err)
		return Guild{}
	}

	Guilds[i.GuildID] = guild

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server added",
		},
	})
	if err != nil {
		log.Println(err)
		return Guild{}
	}

	return guild
}

func removeServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate) Guild {
	guild, done := foundGuild(s, i)
	if done {
		return Guild{}
	}

	appOption := i.ApplicationCommandData().Options[0]
	ip, port := appOption.Options[0].StringValue(), appOption.Options[1].StringValue()

	if !guild.HasServer(ip, port) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Server don't exist",
			},
		})
		if err != nil {
			log.Println(err)
		}
		return Guild{}
	}

	if _, err := guild.RemoveServer(sql.Server{IP: ip, Port: port}); err != nil { // Remove from database
		log.Println(err)
		return Guild{}
	}

	Guilds[i.GuildID] = guild

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server removed",
		},
	})
	if err != nil {
		log.Println(err)
		return Guild{}
	}

	return guild
}

func RefreshServerInfo() {
	for {
		ServerInfo()
		time.Sleep(1 * time.Minute)
	}
}

func ServerInfo() {
	for _, guild := range Guilds {
		DisplayServerInfo(guild)
	}
}

func DisplayServerInfo(guild Guild) {
	mes, err := DG.ChannelMessage(guild.chanelID, guild.messageID)
	if err != nil {
		fmt.Println("error while fetching message: ", err)
		return
	}

	messageUpdate := guild.infos.UpdateMessage()
	messageUpdate.ID = mes.ID
	messageUpdate.Channel = mes.ChannelID

	_, err = DG.ChannelMessageEditComplex(messageUpdate)
	if err != nil {
		fmt.Println(err)
	}
}

func foundGuild(s *discordgo.Session, i *discordgo.InteractionCreate) (Guild, bool) {
	guild, ok := Guilds[i.GuildID]
	if !ok {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Guild not found",
			},
		})
		if err != nil {
			log.Println(err)
		}
		return Guild{}, true
	}
	return guild, false
}
