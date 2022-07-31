package app

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
