package app

import (
	"encoding/json"
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
	header     http.Header
)

type HTTPClient struct {
	driver *http.Client
}

type HTTPError struct {
	StatusCode int
	Status     string
}

func InitClient() {
	// Create HTTP tranport
	transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 15 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 30 * time.Second,
	}

	// List required header
	header = http.Header{
		"Authorization": {"Mzc3MTIzMzk5OTcxMTEwOTEz.GlX-Jv.LIP8_wJ3q_kAxboGyhURMSRHI8Dxt8d-VZFgCM"},
	}

	// Create logfile
	LogFilePtr, err := os.Create(Config.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(os.Stdout, LogFilePtr)
	fileLogger = log.New(w, "", log.LstdFlags)
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		driver: &http.Client{
			Transport: transport,
		},
	}
}

func (client *HTTPClient) GetResponse(url string) (*http.Response, error) {
	// Get current page
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := client.driver.Do(req)

	// Handle status code error
	if resp != nil && resp.StatusCode != 200 {
		err = &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}
	return resp, err
}

func (client *HTTPClient) GetJson(url string, target interface{}) error {
	// Get current page
	resp, err := client.GetResponse(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func (client *HTTPClient) DownloadFile(url string, filepath string) {

	// fmt.Println("Downloading", filepath)

	// Get HTTP response
	resp, err := client.GetResponse(url)
	if err != nil {
		logToFile(err, filepath, url)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		logToFile(err, filepath, url)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logToFile(err, filepath, url)
		return
	}

	// fmt.Println("Download to", filepath, "is finished")
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("Error: %s", e.Status)
}

func logToFile(err error, filepath string, url string) {
	msg := fmt.Sprintf("%v - %s - %s", err, filepath, url)
	fileLogger.Println(msg)
}
