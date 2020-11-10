package discordloghandler

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Handler interface {
	HandleAuth(timestamp, daemon, message string) error
}

type DConfig struct {
	Token    string
	Channels []string
	Users    []string
}

type discordLogHandler struct {
	session  *discordgo.Session
	settings DConfig
}

func (s *discordLogHandler) logSSHDActivity(channel, daemon, username, action, ip, timestamp string) {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "**Username**",
			Value:  username,
			Inline: false,
		},
		{
			Name:   "**IP**",
			Value:  ip,
			Inline: false,
		},
	}

	embed := discordgo.MessageEmbed{
		URL:         "",
		Type:        "rich",
		Title:       "Log",
		Description: action,
		Color:       0,
		Footer: &discordgo.MessageEmbedFooter{
			Text:         daemon,
			IconURL:      "",
			ProxyIconURL: "",
		},
		Image:     nil,
		Thumbnail: nil,
		Video:     nil,
		Provider:  nil,
		Author:    nil,
		Fields:    fields,
		Timestamp: timestamp,
	}
	_, err := s.session.ChannelMessageSendEmbed(channel, &embed)
	if err != nil {
		fmt.Println("Failed to send Embed Message")
	}
}

func (s *discordLogHandler) authLogGeneral(channel, timestamp, daemon, message string) {
	embed := discordgo.MessageEmbed{
		URL:         "",
		Type:        "rich",
		Title:       "Auth Log Message",
		Description: message,
		Color:       0,
		Footer: &discordgo.MessageEmbedFooter{
			Text:         daemon,
			IconURL:      "",
			ProxyIconURL: "",
		},
		Image:     nil,
		Thumbnail: nil,
		Video:     nil,
		Provider:  nil,
		Author:    nil,
		Fields:    nil,
		Timestamp: timestamp,
	}
	_, err := s.session.ChannelMessageSendEmbed(channel, &embed)
	if err != nil {
		fmt.Println("Failed to send Embed Message")
	}
}

func (s *discordLogHandler) HandleAuth(timestamp, daemon, message string) error {
	if s.session == nil {
		return errors.New("discord session isn't initialized")
	}
	for i := 0; i < len(s.settings.Channels); i++ {
		go s.authLogGeneral(s.settings.Channels[i], timestamp, daemon, message)
	}
	for i := 0; i < len(s.settings.Users); i++ {
		channel, err := s.session.UserChannelCreate(s.settings.Users[i])
		if err != nil {
			fmt.Println("Failed to send to user with id: " + s.settings.Users[i])
			continue
		}
		go s.authLogGeneral(channel.ID, timestamp, daemon, message)
	}
	return nil
}

func New(config DConfig) (*discordLogHandler, error) {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	d := &discordLogHandler{
		session:  session,
		settings: config,
	}
	return d, nil
}

func (s *discordLogHandler) Open() error {
	if s.session == nil {
		return errors.New("discord session isn't initialized")
	}
	return s.session.Open()
}

func (s *discordLogHandler) Close() error {
	if s.session == nil {
		return errors.New("discord session isn't initialized")
	}
	return s.session.Close()
}
