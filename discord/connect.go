package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var DG *discordgo.Session

// NewAuth initiates a new discordgo session with the provided Discord token.
// The token is fetched from the "DISCORD_TOKEN" environment variable.
// If there is an error while creating or opening the session, the function
// will terminate the program and log the encountered error.
// Returns a discordgo session object if the process is successful.
func NewAuth() {
	// Replace "YOUR_TOKEN" with your Discord token
	token := os.Getenv("DISCORD_TOKEN")
	var err error

	// Create a Discord session
	DG, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session :", err)
	}

	// Attempt to open the Discord session
	err = DG.Open()
	if err != nil {
		log.Fatal("Error opening Discord session :", err)
	}
}
