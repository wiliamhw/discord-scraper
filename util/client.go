package util

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	LogFilePtr *os.File
	fileLogger *log.Logger
	transport  *http.Transport
)

type HTTPClient struct {
	driver *http.Client
}

type HTTPError struct {
	StatusCode int
	Status     string
}

func init() {
	// Create HTTP tranport
	transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 15 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 30 * time.Second,
	}

	// Create logfile
	LogFilePtr, err := os.Create(Config.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, LogFilePtr)
	fileLogger = log.New(w, "", log.LstdFlags)
}

func (client *HTTPClient) GetResponse(url string) (*http.Response, error) {
	// Get current page
	resp, err := client.driver.Get(url)

	// Handle status code error
	if resp != nil && resp.StatusCode != 200 {
		err = &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}
	return resp, err
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		driver: &http.Client{
			Transport: transport,
		},
	}
}

func (client *HTTPClient) DownloadFile(url string, filepath string) {

	// fmt.Println("Downloading", filepath)

	// Get HTTP response
	resp, err := client.GetResponse(url)
	if err != nil {
		fileLogger.Println(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fileLogger.Println(err)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fileLogger.Println(err)
		return
	}

	// fmt.Println("Download to", filepath, "is finished")
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Status code error: %d %s", e.StatusCode, e.Status)
}
