package discord

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"bot-serveur-info/internal/pkg/controller"
	"bot-serveur-info/pkg/discord"
)

var (
	componentsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild{
		"update": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			if err := discord.BasicResponse(s, i, "Refreshed"); err != nil {
				log.Println(err)
			}
			return guild
		},
		"right": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			guild.NextPage()

			if err := discord.BasicResponse(s, i, "Page "+strconv.Itoa(guild.Infos().Page)); err != nil {
				log.Println(err)
			}
			return guild
		},
		"left": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			guild.PreviousPage()

			if err := discord.BasicResponse(s, i, "Page "+strconv.Itoa(guild.Infos().Page)); err != nil {
				log.Println(err)
			}
			return guild
		},
	}
	commandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild{
		"server add": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			return addServerCommand(s, i, guild)
		},
		"server remove": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			return removeServerCommand(s, i, guild)
		},
		"server message": func(s *discordgo.Session, i *discordgo.InteractionCreate, guild controller.Guild) controller.Guild {
			return setMessage(s, i, guild)
		},
	}
)
