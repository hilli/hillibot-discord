package bot

import (
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hilli/hillibot-discord/config"
	"github.com/rs/zerolog/log"
)

var (
	BotID string
	goBot *discordgo.Session
)

func Start() {
	log.Info().Msg("Connecting to Discord")
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Error().Err(err).Msg("Could not connect to Discord")
		os.Exit(1)
	}
	u, err := goBot.User("@me")
	if err != nil {
		log.Error().Err(err).Msg("Could get Discord user")
		os.Exit(2)
	}

	BotID = u.ID
	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		log.Error().Err(err).Msg("Could not connect to Discord websocket")
		os.Exit(2)
		return
	}
	log.Info().Msg("Bot is running")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't respond to myself
	if m.Author.ID == BotID {
		return
	}

	// ping:pong
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
		log.Info().Msgf("%s said ping to us", m.Author.Username)
	}

	// Finance graph
	if strings.HasPrefix(m.Content, "!fin") {
		financeHandler(s, m)
		log.Info().Msgf("%s said '%s' to us", m.Author.Username, m.Content)
	}

}
