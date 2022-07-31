package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-module/carbon"
	"github.com/wiliamhw/discord-scraper/app"
	"github.com/wiliamhw/discord-scraper/model"
	"github.com/wiliamhw/discord-scraper/util"
)

const (
	baseURL  = "https://discord.com/api/v9"
	basePath = "storage/results"
)

func main() {
	start := time.Now()
	app.InitConfig()
	fmt.Printf("Downloading Discord channel: %s\n", app.Input.ChannelId)
	app.InitClient()
	app.InitWorker()
	defer app.LogFilePtr.Close()

	// Get storage path
	now := time.Now().Format("2006-01-02 15:04:05")
	nowFormatted := carbon.Parse(now).Format("Y-m-d_H-i-s")
	storagePath := fmt.Sprintf("%s/%s", basePath, nowFormatted)

	// Get all chats in a channel
	client := app.NewHTTPClient()
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
		baseURL, app.Input.ChannelId, app.Input.NumOfChats,
	)
	var chats []model.Chat
	err := client.GetJson(url, &chats)
	if err != nil {
		log.Fatal(err)
	}

	// Download every attachments in each chat
	for index, chat := range chats {
		if len(chat.Attachments) == 0 {
			continue
		}
		timestamp := carbon.Parse(chat.Timestamp).Format("Y-m-d_H-i-s")

		// Donwload as a file if only 1 attachment exists
		if len(chat.Attachments) == 1 {
			attachment := chat.Attachments[0]
			ext := filepath.Ext(attachment.URL)
			filePath := fmt.Sprintf("%s/%s%s", storagePath, timestamp, ext)
			attachment.Download(index, filePath)
			continue
		}

		// Create folder to download >1 attachments in a folder
		dirPath := fmt.Sprintf("%s/%s", storagePath, timestamp)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		// Download file
		for _, attachment := range chat.Attachments {
			filePath := fmt.Sprintf("%s/%s", dirPath, attachment.Filename)
			attachment.Download(index, filePath)
		}
	}

	close(app.JobsQueue)
	app.WG.Wait()
	util.PruneEmptyDirectories(storagePath)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
