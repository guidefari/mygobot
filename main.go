package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	outfile, _            = os.Create("server.log")
	local_log             = log.New(outfile, "", log.LstdFlags|log.Lshortfile)
	discord_session_token string
	preferred_channel_id  string
	disco_session         *discordgo.Session
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	env_var, exists := os.LookupEnv("DISCORD_TOKEN")
	if exists {
		fmt.Println("Discord token found & loaded")
		discord_session_token = env_var
	}

	env_var, exists = os.LookupEnv("PREFERRED_CHANNEL_ID")
	if exists {
		fmt.Println("Channel ID we in")
		preferred_channel_id = env_var
	}

	sess, err := discordgo.New("Bot " + discord_session_token)
	if err != nil {
		log.Fatal(err)
	}

	disco_session = sess
	err = disco_session.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is online")
}

func main() {

	disco_session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	disco_session.AddHandler(session_handler)
	defer disco_session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func session_handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	messageContent := m.Content

	if messageContent == "!rss" {
		s.ChannelMessageSend(m.ChannelID, "Yo! I'm here. Use `!rss help`")
		local_log.Printf("DCScragor got a command from %s, Command: %s", m.Author, messageContent)
	}

	// !rss help
	if strings.HasPrefix(m.Content, "!rss") && len(strings.Split(m.Content, " ")) >= 2 {
		if strings.Split(m.Content, " ")[1] == "help" {
			sendDudes(s, m)
		}
	}

}

func sendDudes(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Available commands:\n----------------------------\n!rss add_blog <URL>       Adds your RSS feed URL to the scrape wish list.")
}
