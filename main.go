package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wiliamhw/discord-scraper/app"
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

	// Catch Ctrl+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(app.JobsQueue)
		os.Exit(1)
	}()

	// Scrape
	storagePath := fmt.Sprintf("%s/%s", basePath, app.Input.DirName)
	scraper := app.RunAPI
	if app.Input.UseJSON {
		scraper = app.RunJSON
	}
	err := scraper(baseURL, storagePath)
	if err != nil {
		log.Fatal(err)
	}

	close(app.JobsQueue)
	app.WG.Wait()
	util.PruneEmptyDirectories(storagePath)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
}
