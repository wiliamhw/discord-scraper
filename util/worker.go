package util

import (
	"fmt"
	"sync"
)

var (
	jobsQueue chan downloadJob
	clients   []*HTTPClient
	WG        = sync.WaitGroup{}
)

type downloadJob struct {
	function     func(interface{})
	arguments    []string
	startMessage string
	endMessage   string
}

func worker(id int, jobs <-chan downloadJob) {
	for j := range jobs {
		fmt.Printf("Worker %02d - %s\n", id, j.startMessage)
		j.function(j.arguments)
		WG.Done()
		fmt.Printf("Worker %02d - Download to %s is finished\n", id, j.endMessage)
	}
}
