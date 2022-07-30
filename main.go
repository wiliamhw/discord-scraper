package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-module/carbon"
	"github.com/wiliamhw/discord-scraper/app"
	"github.com/wiliamhw/discord-scraper/util"
)

const baseURL = "https://discord.com/api/v9"
const storagePath = "storage"

type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

type Chat struct {
	ID          string       `json:"id"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   string       `json:"timestamp"`
}

func main() {
	start := time.Now()
	app.InitConfig()
	fmt.Printf("Downloading Discord channel: %s\n", app.Input.ChannelId)
	app.InitClient()
	app.InitWorker()
	defer app.LogFilePtr.Close()

	// Get all chats in a channel
	client := app.NewHTTPClient()
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
		baseURL, app.Input.ChannelId, app.Input.NumOfChats,
	)
	var chats []Chat
	err := client.GetJson(url, &chats)
	if err != nil {
		log.Fatal(err)
	}

	// Download every attachments in each chat
	for index, chat := range chats {
		if len(chat.Attachments) == 0 {
			continue
		}

		// Create folder
		timestamp := carbon.Parse(chat.Timestamp).Format("Y-m-d_H-i-s")
		dirPath := fmt.Sprintf("%s/%s", storagePath, timestamp)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		// Download file
		for _, attachment := range chat.Attachments {
			filePath := fmt.Sprintf("%s/%s", dirPath, attachment.Filename)
			app.WG.Add(1)
			app.JobsQueue <- app.DownloadJob{
				Number: index,
				URL:    attachment.URL,
				Path:   filePath,
			}
		}
	}

	close(app.JobsQueue)
	app.WG.Wait()
	util.PruneEmptyDirectories(storagePath)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
