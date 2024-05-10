package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var discord_session_token string
var preferred_channel_id string
var disco_session *discordgo.Session

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	token, exists := os.LookupEnv("DISCORD_TOKEN")
	if exists {
		fmt.Println("Discord token found & loaded")
		discord_session_token = token
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

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "Thanks for the refactor🫡")
	}

}
