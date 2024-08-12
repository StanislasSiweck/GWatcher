package session

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/lmittmann/tint"

	"bot-serveur-info/pkg/logger"
)

var DG *discordgo.Session

func NewAuth() {
	// Replace "YOUR_TOKEN" with your Discord token
	token := os.Getenv("DISCORD_TOKEN")
	var err error

	// Create a Discord session
	DG, err = discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatal("Can't creating Discord session", tint.Err(err))
	}

	// Attempt to open the Discord session
	err = DG.Open()
	if err != nil {
		logger.Fatal("Can't opening Discord session", tint.Err(err))
	}
}
