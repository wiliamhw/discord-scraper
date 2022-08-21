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

func exit(start time.Time, storagePath string, exitCode int) {
	close(app.JobsQueue)
	if exitCode == 0 {
		app.WG.Wait()
	}
	util.PruneEmptyDirectories(storagePath)
	fmt.Printf("\nTook: %f secs\n", time.Since(start).Seconds())
	app.LogFilePtr.Close()
	os.Exit(exitCode)
}

func main() {
	start := time.Now()
	app.InitConfig()
	fmt.Printf("Downloading Discord channel: %s\n", app.Input.ChannelId)
	app.InitClient()
	app.InitWorker()

	// Create folder
	storagePath := fmt.Sprintf("%s/%s", basePath, app.Input.DirName)
	dirPath := fmt.Sprintf("%s", storagePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Catch Ctrl+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exit(start, storagePath, 1)
	}()

	// Scrape
	scraper := app.RunAPI
	if app.Input.UseJSON {
		scraper = app.RunJSON
	}
	err := scraper(baseURL, storagePath)
	if err != nil {
		log.Fatal(err)
	}

	exit(start, storagePath, 0)
}
