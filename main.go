package main

import (
	"fmt"
	"time"

	"github.com/wiliamhw/discord-scraper/util"
)

func main() {
	start := time.Now()
	defer util.LogFilePtr.Close()
	fmt.Printf("Downloading Discord channel: %s\n", util.Input.ChannelId)

	close(util.JobsQueue)
	util.WG.Wait()
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
