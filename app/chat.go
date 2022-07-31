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
	client := NewHTTPClient()
	url := fmt.Sprintf("%s/channels/%s/messages?limit=%d",
		baseURL, Input.ChannelId, Input.NumOfChats,
	)
	var chats []Chat
	if err := client.GetJson(url, &chats); err != nil {
		return nil, err
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
