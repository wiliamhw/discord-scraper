package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang-module/carbon"
)

// Download every attachments in each chat
func parseChats(chats *[]Chat, storagePath string) error {
	// Download every attachments in each chat
	for index, chat := range *chats {
		if index >= Input.NumOfChats {
			break
		}

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
			return err
		}

		// Download file
		for _, attachment := range chat.Attachments {
			filePath := fmt.Sprintf("%s/%s", dirPath, attachment.Filename)
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
