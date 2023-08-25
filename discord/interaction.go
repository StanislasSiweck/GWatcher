package discord

import (
	"bot-serveur-info/sql"
	"github.com/bwmarrin/discordgo"
	"log"
)

var AllServers = map[string]sql.Server{}

// InteractionCreate is a function that handles when a new Discord interaction is created.
//
// Parameters:
// s *discordgo.Session: The current Discord session.
// i *discordgo.InteractionCreate: Represents the creation of a new interaction in Discord.
//
// The function checks if the interaction is of type ApplicationCommand.
// If it's not, the function immediately returns.
// If it is, the function retrieves the command data from the interaction and
// identifies the command name. Depending on the command name, it calls the respective function.
// For 'addserver', it calls addServerCommand and for 'removeserver', it calls removeServerCommand.
func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	command := data.Name
	switch command {
	case "addserver":
		addServerCommand(s, i, data)
	case "removeserver":
		removeServerCommand(s, i, data)
	}
}

// addServerCommand is a function that handles the command to add a new server
// to the database and a local server list.
//
// Parameters:
// s *discordgo.Session - The current Discord session.
// i *discordgo.InteractionCreate - Represents an interaction creation event from Discord.
// data discordgo.ApplicationCommandInteractionData - Represents the data of an application command interaction.
//
// The function retrieves the server details (game, IP, and port) from the data options,
// and creates a new server instance with those details.
// If the creation in the database is successful, it also adds the server to the AllServers map.
// Then it sends a response to the Discord 'interaction' confirming that the server has been added.
// If there are any errors along the way, these are logged.
func addServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	server := sql.Server{
		IP:   data.Options[0].StringValue(),
		Port: data.Options[1].StringValue(),
	}

	if err := sql.AddServer(server); err != nil { // Create the server in the database
		log.Println(err)
	}

	AllServers[server.IP+":"+server.Port] = server // Add to local list

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server added",
		},
	})
	if err != nil {
		log.Println(err)
	}
}

// removeServerCommand is a function that handles the removal of a server
// from the local server list and the database.
//
// Parameters:
// s *discordgo.Session - The current Discord session.
// i *discordgo.InteractionCreate - Represents an interaction creation event from Discord.
// data discordgo.ApplicationCommandInteractionData - Represents the data of an application command interaction.
//
// The function retrieves the server details (IP, and port) from the data options,
// and then attempts to remove the server from the AllServers map and the database.
// Upon successful removal of the server, it will send a message back to the Discord channel
// through the provided session and interaction, notifying that the server has been removed.
// If there are any errors along the way, they are logged.
func removeServerCommand(s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	delete(AllServers, data.Options[0].StringValue()+":"+data.Options[1].StringValue()) // Remove from local list

	ip, port := data.Options[0].StringValue(), data.Options[1].StringValue()

	if err := sql.RemoveServer(ip, port); err != nil { // Remove from database
		log.Println(err)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{ // Send response to Discord
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Server removed",
		},
	})
	if err != nil {
		log.Println(err)
	}
}
