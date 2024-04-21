package discord

import (
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild{
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			if err := discord.BasicResponse(s, i, "Refreshed"); err != nil {
				log.Println(err)
			}
			return guild
		},
		"right": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			guild.NextPage()

			if err := discord.BasicResponse(s, i, "Page "+strconv.Itoa(guild.Infos.Page)); err != nil {
				log.Println(err)
			}
			return guild
		},
		"left": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			guild.PreviousPage()

			if err := discord.BasicResponse(s, i, "Page "+strconv.Itoa(guild.Infos.Page)); err != nil {
				log.Println(err)
			}
			return guild
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild{
		"server add": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			return addServerCommand(s, i, guild)
		},
		"server remove": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			return removeServerCommand(s, i, guild)
		},
		"server message": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
			return setMessage(s, i, guild)
		},
	}
)
