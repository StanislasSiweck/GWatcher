package discord

import (
	"bot-serveur-info/internal/pkg/class"
	"bot-serveur-info/internal/pkg/session"
	"bot-serveur-info/internal/pkg/sql/model"
	"bot-serveur-info/internal/pkg/sql/request"
	"bot-serveur-info/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"
	"log/slog"
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
		slog.Error("Can't convert guild id to int", tint.Err(err), "guild_id", g.ID)
		return
	}
	guild, err := class.CreateGuild(uint(guildID), "", "")
	if err != nil {
		slog.Error("Can't create guild", tint.Err(err), "guild_id", g.ID)
		return
	}
	guild.SetDisplayInfo(class.NewDisplay([]model.Server{}, 0))
	Guilds[g.ID] = guild
}

func GuildDelete(_ *discordgo.Session, g *discordgo.GuildDelete) {
	_, ok := Guilds[g.ID]
	if !ok {
		return
	}

	if err := request.RemoveGuild(g.ID); err != nil {
		slog.Error("Can't remove guild", tint.Err(err), "guild_id", g.ID)
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
				slog.Error("Guild not found", "guild_id", i.GuildID)
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
				slog.Error("Guild not found", "guild_id", i.GuildID)
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
		slog.Error("Can't send message", tint.Err(err), "channel_id", i.ChannelID)
		return class.Guild{}
	}

	guild.ChanelID = mes.ChannelID
	guild.MessageID = mes.ID
	if err := guild.UpdateGuild(); err != nil {
		slog.Error("Can't update guild", tint.Err(err), "guild_id", mes.ID)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Message send"); err != nil {
		slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
		return guild
	}

	return guild
}

func addServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
	appOption := i.ApplicationCommandData().Options[0]

	ip, port := appOption.Options[0].StringValue(), appOption.Options[1].StringValue()

	if guild.HasServer(ip, port) {
		if err := discord.BasicResponse(s, i, "Server already existed"); err != nil {
			slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
		}
		return class.Guild{}
	}

	server := model.Server{
		IP:   ip,
		Port: port,
	}

	if _, err := guild.AddServer(server); err != nil { // Create the server in the database
		slog.Error("Can't add server", tint.Err(err), "IP", ip, "Port", port)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Server added"); err != nil {
		slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
		return guild
	}

	return guild
}

func removeServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, guild class.Guild) class.Guild {
	appOptions := i.ApplicationCommandData().Options[0]
	ip, port := appOptions.Options[0].StringValue(), appOptions.Options[1].StringValue()

	if !guild.HasServer(ip, port) {
		if err := discord.BasicResponse(s, i, "Server don't exist"); err != nil {
			slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
		}
		return class.Guild{}
	}

	if _, err := guild.RemoveServer(model.Server{IP: ip, Port: port}); err != nil { // Remove from database
		slog.Error("Can't remove server", tint.Err(err), "IP", ip, "Port", port)
		return class.Guild{}
	}

	if err := discord.BasicResponse(s, i, "Server removed"); err != nil {
		slog.Error("Can't send a basic reply", tint.Err(err), "guild_id", i.GuildID)
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
