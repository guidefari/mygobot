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
	var command []string

	if strings.HasPrefix(messageContent, "!rss") {
		command = strings.Split(messageContent, " ")
	}

	if len(command) == 0 {
		return
	}

	if messageContent == "!rss" {
		s.ChannelMessageSend(m.ChannelID, "Yo! I'm here. Use `!rss help`")
		local_log.Printf("DCScragor got a command from %s, Command: %s", m.Author, messageContent)
	}

	// !rss help
	if len(command) >= 2 {
		if command[1] == "help" {
			sendDudes(s, m)
		}
	}

	// !rss add_blog <url>
	if len(command) >= 3 {
		if command[1] == "add_blog" && isUrl(command[2]) {
			s.ChannelMessageSend(m.ChannelID, "Adding blog to list 🔥")
			handleAddBlogToList(s, command[2], m)
			local_log.Printf("Blog adding request from %s, Command: %s", m.Author, m.Content)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Please use a proper URL!")
		}
	}

}

func sendDudes(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Available commands:\n----------------------------\n!rss add_blog <URL>       Adds your RSS feed URL to the scrape wish list.")
}

func isUrl(s string) bool {
	return strings.Contains(s, "http://") || strings.Contains(s, "https://")
}

func handleAddBlogToList(s *discordgo.Session, url string, m *discordgo.MessageCreate) {
	blog_list_filename := "blog_request.list"

	file, err := os.OpenFile(blog_list_filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil && err.Error() == "open blog_request.list: no such file or directory" {
		file, _ = os.Create(blog_list_filename)
		local_log.Printf("File Created: blog_request.list")
	}

	defer file.Close()
	if !strings.Contains(readFile(blog_list_filename), url) {
		if _, err := file.WriteString(url + "\n"); err != nil {
			local_log.Fatal(err)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Blog URL is added to wish list.")
			local_log.Printf("Added a new blog to the wish list. Author: %s, URL: %s", m.Author, url)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "You've already requested for this URL!")
	}
}
