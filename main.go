package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
	"sync"
	"time"
)

type config struct {
	Token     *string
	Channels  []*string
	Users     []*string
	NamedPipe *string
	Daemons   []*string
}

func logSSHDActivity(session *discordgo.Session, channel, daemon, username, action, ip, timestamp string) {
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
	_, err := session.ChannelMessageSendEmbed(channel, &embed)
	if err != nil {
		fmt.Println("Failed to send Embed Message")
	}
}

func authLogGeneral(session *discordgo.Session, channel, timestamp, daemon, message string) {
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
	_, err := session.ChannelMessageSendEmbed(channel, &embed)
	if err != nil {
		fmt.Println("Failed to send Embed Message")
	}
}

func parseJson(filePath string) (*config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open settings file")
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Failed getting file stat")
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Failed to read file")
		return nil, err
	}
	var settings config
	err = json.Unmarshal(buffer, &settings)
	if err != nil {
		fmt.Println("Failed to unmarshal")
		return nil, err
	}
	return &settings, nil
}

func discordHandler(settings *config, timestamp, daemon, message string) {
	session, err := discordgo.New("Bot " + *settings.Token)
	if err != nil {
		fmt.Println("Failed to start bot, invalid settings")
		return
	}

	err = session.Open()
	if err != nil {
		fmt.Println("Failed to connect to discord")
		return
	}
	defer session.Close()

	for i := 0; i < len(settings.Channels); i++ {
		go authLogGeneral(session, *settings.Channels[i], timestamp, daemon, message)
	}
	for i := 0; i < len(settings.Users); i++ {
		channel, err := session.UserChannelCreate(*settings.Users[i])
		if err != nil {
			fmt.Println("Failed to send to user with id: " + *settings.Users[i])
			continue
		}
		go authLogGeneral(session, channel.ID, timestamp, daemon, message)
	}
}

func stringContainsOneOfSubstringArray(a string, b []*string) bool {
	for _, i := range b {
		if strings.Contains(a, *i) {
			return true
		}
	}
	return false
}

func handleLog(logLine string, settings *config) {
	listOfStrings := strings.SplitN(logLine, " ", 4)
	timestamp := listOfStrings[0]
	daemon := strings.Trim(listOfStrings[2], ":")
	message := listOfStrings[3]
	if stringContainsOneOfSubstringArray(daemon, settings.Daemons) {
		discordHandler(settings, timestamp, daemon, message)
	}
}

func handleNamedPipeSyslog(settings *config, wg *sync.WaitGroup) {
	namedPipe, _ := os.OpenFile(*settings.NamedPipe, os.O_RDONLY, os.ModeNamedPipe)
	defer namedPipe.Close()
	defer wg.Done()
	for {
		scanner := bufio.NewScanner(namedPipe)
		for scanner.Scan() {
			logMessage := scanner.Text()
			logLines := strings.Split(logMessage, "\n")
			for i := 0; i < len(logLines); i++ {
				logLine := logLines[i]
				handleLog(logLine, settings)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func main() {
	settings, err := parseJson("settings.json")
	if err != nil {
		fmt.Println("Failed to parse json")
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go handleNamedPipeSyslog(settings, &wg)
	wg.Wait()
}
