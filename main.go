package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-module/carbon"
	"github.com/wiliamhw/discord-scraper/app"
	"github.com/wiliamhw/discord-scraper/scraper"
	"github.com/wiliamhw/discord-scraper/util"
)

const (
	baseURL  = "https://discord.com/api/v9"
	basePath = "storage/results"
)

func main() {
	start := time.Now()
	app.InitConfig()
	fmt.Printf("Downloading Discord channel: %s\n", app.Input.ChannelId)
	app.InitClient()
	app.InitWorker()
	defer app.LogFilePtr.Close()

	// Get storage path
	now := time.Now().Format("2006-01-02 15:04:05")
	nowFormatted := carbon.Parse(now).Format("Y-m-d_H-i-s")
	storagePath := fmt.Sprintf("%s/%s", basePath, nowFormatted)

	// Scrape
	functionName := scraper.RunJSON
	if app.Input.UseJSON {
		functionName = scraper.RunAPI
	}
	err := functionName(baseURL, storagePath)
	if err != nil {
		log.Fatal(err)
	}

	close(app.JobsQueue)
	app.WG.Wait()
	util.PruneEmptyDirectories(storagePath)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
