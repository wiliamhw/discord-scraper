package scraper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-module/carbon"
	"github.com/wiliamhw/discord-scraper/app"
)

func RunAPI(baseURL string, storagePath string) (err error) {
	// Get all chats in a channel
	client := app.NewHTTPClient()
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
		baseURL, app.Input.ChannelId, app.Input.NumOfChats,
	)
	var chats []app.Chat
	err = client.GetJson(url, &chats)
	if err != nil {
		return
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
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return
		}

		// Download file
		for _, attachment := range chat.Attachments {
			filePath := fmt.Sprintf("%s/%s", dirPath, attachment.Filename)
			attachment.Download(index, filePath)
		}
	}
	return
}
