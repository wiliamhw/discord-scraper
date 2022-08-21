package app

import (
	"encoding/json"
	"fmt"
	"os"
)

type Chat struct {
	ID          string       `json:"id"`
	Content     string       `json:"content"`
	Attachments []Attachment `json:"attachments"`
	Timestamp   string       `json:"timestamp"`
}

type Attachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
}

func (a *Attachment) Download(number int, filePath string) {
	WG.Add(1)
	JobsQueue <- DownloadJob{
		Number: number,
		URL:    a.URL,
		Path:   filePath,
	}
}

////////////////////////////////
// Get all chats in a channel
////////////////////////////////

func GetChatsFromAPI(baseURL string) (*[]Chat, error) {
	var chats []Chat
	totalChat := 0

	client := NewHTTPClient()
	client.WithHeader = true

	for totalChat < Input.NumOfChats {
		// Determine limit paramms
		var currentChats []Chat
		limit := 100
		if Input.NumOfChats-totalChat < 100 {
			limit = Input.NumOfChats - totalChat
		}

		// Determine URI to fetch
		url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
			baseURL, Input.ChannelId, limit,
		)
		lenChats := len(chats)
		if lenChats > 0 {
			url += "?before=" + chats[lenChats-1].ID
		}

		// Fetch API
		if err := client.GetJson(url, &currentChats); err != nil {
			return nil, err
		}
		chats = append(chats, currentChats...)
		totalChat += limit
	}

	return &chats, nil
}

func GetChatsFromJSON(sourcePath string) (*[]Chat, error) {
	// Open file
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}

	// Parse JSON in the opened file
	dec := json.NewDecoder(file)
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if t != json.Delim('{') {
		return nil, fmt.Errorf("Expected {, got %v", t)
	}
	for dec.More() {

		// Read the key.
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key := t.(string) // type assert token to string.

		// Return messages array
		if key == "messages" {
			var value []Chat
			if err := dec.Decode(&value); err != nil {
				return nil, err
			}
			return &value, nil
		}

		// Decode the value.
		var value interface{}
		if err := dec.Decode(&value); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("Field \"messages\" not found")
}
