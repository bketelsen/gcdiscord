package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token      string
	GuildRoles []*discordgo.Role
)

const GopherGuild = "755435423177638059"

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(memberJoin)
	dg.AddHandler(presenceUpdate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	GuildRoles, err = dg.GuildRoles(GopherGuild)
	if err != nil {
		fmt.Println("Error getting roles:", err)
	}
	for _, gr := range GuildRoles {
		describeRole(gr.ID)
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func presenceUpdate(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	fmt.Println("Presence Update!")
	ensureRoles(p)
}
func memberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {

	fmt.Println("Member Join!")
	fmt.Println(m.Member)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message Create!")
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func ensureRoles(p *discordgo.PresenceUpdate) error {
	fmt.Println("Ensure Roles...")
	fmt.Println("User ID:", p.User.ID)
	fmt.Println("Username:", p.User.Username)
	for _, r := range p.Roles {
		fmt.Println("Role ID:", r)
		describeRole(r)
	}
	return nil
}

func describeRole(r string) {
	for _, role := range GuildRoles {
		if role.ID == r {
			fmt.Println("Role Name: ", role.Name)
		}
	}
}
