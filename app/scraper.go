package app

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/golang-module/carbon"
)

// Download every attachments in each chat
func parseChats(chats *[]Chat, storagePath string) error {
	fmt.Printf("Parsing %d chats\n", len(*chats))

	// Download every attachments in each chat
	for index, chat := range *chats {
		if Input.NumOfChats > 0 && index >= Input.NumOfChats {
			break
		}

		if len(chat.Attachments) == 0 {
			continue
		}
		timestamp := carbon.Parse(chat.Timestamp).Format("Y-m-d_H-i-s")
		if timestamp == "" {
			timestamp = carbon.Now().Format("Y-m-d_H-i-s")
		}

		// Download file
		for idx, attachment := range chat.Attachments {
			ext := filepath.Ext(attachment.URL)
			filePath := fmt.Sprintf(
				"%s/%s_%d%s",
				storagePath, timestamp, idx+1, ext,
			)
			attachment.Download(index, filePath)
		}
	}
	return nil
}

func RunAPI(baseURL string, storagePath string) error {
	chats, err := GetChatsFromAPI(baseURL)
	if err != nil {
		return err
	}

	return parseChats(chats, storagePath)
}

func RunJSON(baseURL string, storagePath string) error {
	sourcePath := "storage/sources"

	// Read all flles in the directory
	fileMetas, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		return err
	}

	// Iterate each flles in the directory
	for _, fileMeta := range fileMetas {
		if filepath.Ext(fileMeta.Name()) != ".json" {
			continue
		}
		fmt.Println("Downloading: ", fileMeta.Name())

		// Get all chats from file
		filePath := fmt.Sprintf("%s/%s", sourcePath, fileMeta.Name())
		chats, err := GetChatsFromJSON(filePath)
		if err != nil {
			return err
		}

		return parseChats(chats, storagePath)
	}
	return nil
}
