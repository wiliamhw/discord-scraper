package util

import (
	"fmt"
	"sync"
)

var (
	JobsQueue chan downloadJob
	WG        = sync.WaitGroup{}
	clients   []*HTTPClient
)

type downloadJob struct {
	URL          string
	destFullPath string
}

func init() {
	clients = make([]*HTTPClient, Config.NumOfWorkers)
	JobsQueue = make(chan downloadJob, Config.JobsBuffer)

	for i := 0; i < Config.NumOfWorkers; i++ {
		client := NewHTTPClient()
		clients = append(clients, client)
		go worker(i+1, client, JobsQueue)
	}
}

func worker(id int, client *HTTPClient, jobs <-chan downloadJob) {
	for j := range jobs {
		fmt.Printf("Worker %02d - Downloading %s\n", id, j.destFullPath)
		client.DownloadFile(j.URL, j.destFullPath)
		WG.Done()
		fmt.Printf("Worker %02d - Download to %s is finished\n", id, j.destFullPath)
	}
}
