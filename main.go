package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wiliamhw/discord-scraper/util"
)

const baseURL = "https://discord.com/api/v9"

type Chat struct {
	Content     string                   `json:"content"`
	Attachments []map[string]interface{} `json:"attachments"`
	Timestamp   string                   `json:"timestamp"`
}

func main() {
	start := time.Now()
	util.InitConfig()
	util.InitClient()
	defer util.LogFilePtr.Close()
	fmt.Printf("Downloading Discord channel: %s\n", util.Input.ChannelId)

	// Get all chats in a channel
	client := util.NewHTTPClient()
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
		baseURL, util.Input.ChannelId, util.Input.NumOfChats,
	)
	var chats []Chat
	err := client.GetJson(url, &chats)
	if err != nil {
		log.Fatal(err)
	}

	close(util.JobsQueue)
	util.WG.Wait()
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
