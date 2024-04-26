package session

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var DG *discordgo.Session

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
