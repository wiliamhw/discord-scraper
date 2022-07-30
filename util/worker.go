package util

import (
	"fmt"
	"sync"
)

var (
	JobsQueue chan DownloadJob
	WG        = sync.WaitGroup{}
	clients   []*HTTPClient
)

type DownloadJob struct {
	Number int
	URL    string
	Path   string
}

func InitWorker() {
	clients = make([]*HTTPClient, Config.NumOfWorkers)
	JobsQueue = make(chan DownloadJob, Config.JobsBuffer)

	for i := 0; i < Config.NumOfWorkers; i++ {
		client := NewHTTPClient()
		clients = append(clients, client)
		go worker(i+1, client, JobsQueue)
	}
}

func worker(id int, client *HTTPClient, jobs <-chan DownloadJob) {
	for j := range jobs {
		fmt.Printf("Worker %02d - %03d - Downloading %s\n",
			id, j.Number, j.Path,
		)
		client.DownloadFile(j.URL, j.Path)
		WG.Done()
		fmt.Printf("Worker %02d - %03d - Download to %s is finished\n",
			id, j.Number, j.Path,
		)
	}
}
