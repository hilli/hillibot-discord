package bot

import (
	"bytes"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hilli/hillibot-discord/finance"
)

func financeHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	symlist := strings.SplitN(m.Content, "!fin ", 2)[1]
	G := bytes.NewBuffer([]byte{})
	finance.RenderedGraphForSymbols(finance.SymbolsToList(symlist), G)
	// Create embedded for graph file
	file := discordgo.File{
		ContentType: "image/png",
		Name:        symlist + time.Now().String() +".png",
		Reader:      bytes.NewReader(G.Bytes()),
	}
	// Upload "complex" message
	s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: "Graph for " + symlist,
		Files:   []*discordgo.File{&file},
	})
}
