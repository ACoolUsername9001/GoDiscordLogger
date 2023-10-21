package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"acooldomain.co/discordLogger/discordloghandler"
	"github.com/mitchellh/mapstructure"
)

func parseJson(filePath string) (map[string]interface{}, error) {
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
	var settings map[string]interface{}
	err = json.Unmarshal(buffer, &settings)
	if err != nil {
		fmt.Println("Failed to unmarshal")
		return nil, err
	}
	return settings, nil
}

func stringContainsOneOfSubstringArray(a string, b []interface{}) bool {
	for _, i := range b {
		if strings.Contains(a, i.(string)) {
			return true
		}
	}
	return false
}

func handleLog(handler discordloghandler.Handler, logLine string, settings map[string]interface{}) {
	listOfStrings := strings.SplitN(logLine, " ", 4)
	timestamp := listOfStrings[0]
	daemon := strings.Trim(listOfStrings[2], ":")
	message := listOfStrings[3]
	if stringContainsOneOfSubstringArray(daemon, settings["daemons"].([]interface{})) {
		_ = handler.HandleAuth(timestamp, daemon, message)
	}
}

func handleNamedPipeSyslog(handler discordloghandler.Handler, settings map[string]interface{}, wg *sync.WaitGroup) {
	namedPipePath := settings["namedPipe"].(string)
	namedPipe, _ := os.OpenFile(namedPipePath, os.O_RDONLY, os.ModeNamedPipe)
	defer namedPipe.Close()
	defer wg.Done()
	for {
		scanner := bufio.NewScanner(namedPipe)
		for scanner.Scan() {
			logMessage := scanner.Text()
			logLines := strings.Split(logMessage, "\n")
			for i := 0; i < len(logLines); i++ {
				logLine := logLines[i]
				handleLog(handler, logLine, settings)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	settings, err := parseJson("settings.json")
	if err != nil {
		fmt.Println("Failed to parse json")
		return
	}
	discordSettingsDict, _ := settings["discord"].(map[string]interface{})
	var discordSettings discordloghandler.DConfig
	err = mapstructure.Decode(&discordSettingsDict, &discordSettings)
	if err != nil {
		fmt.Println("Failed to load discord config")
		return
	}
	discordHandler, err := discordloghandler.New(discordSettings)
	if err != nil {
		fmt.Println("Failed to create discordHandler")
		return
	}
	err = discordHandler.Open()
	if err == nil {
		defer discordHandler.Close()
	} else {
		fmt.Println("Failed to load discord")
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go handleNamedPipeSyslog(discordHandler, settings, &wg)
	wg.Wait()
}
