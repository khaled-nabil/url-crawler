package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func main() {
	loadEnv()

	urlsEnv := os.Getenv("URLS")
	folderPath := os.Getenv("FOLDER_PATH")
	urls := strings.Split(urlsEnv, ",")

	schedule := time.Now().Add(time.Hour * 24 * 7)
	timer := time.NewTimer(schedule.Sub(time.Now()))
	fmt.Println("starting the crawler...", urls, time.Now())

	// make folder folderPath if doesn't exist
	err := os.Mkdir(folderPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating parent folder:", err)
		return
	}

	// first unscheduled crawl
	crawl(urls, folderPath)

	for {
		<-timer.C

		crawl(urls, folderPath)

		timer.Reset(time.Hour * 24 * 7)
	}
}

func crawl(urls []string, folderPath string) {
	fmt.Println("crawling in my skin", time.Now())

	folderName := time.Now().Format("02.01.06")
	outputFolderPath := fmt.Sprintf("%s/%s", folderPath, folderName)

	err := os.Mkdir(outputFolderPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating folder:", err)
		return
	}

	// Process each URL
	for _, url := range urls {
		// Extract domain name from URL
		domain := extractDomain(url)

		// Make a GET request
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error making HTTP request:", err, url)
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		// Save the response body to a file
		fileName := fmt.Sprintf("%s/%s.txt", outputFolderPath, domain)
		err = ioutil.WriteFile(fileName, body, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		fmt.Printf("Saved %s to %s\n", url, fileName)
	}
}

func extractDomain(url string) string {
	parts := strings.Split(url, "//")
	if len(parts) > 1 {
		url = parts[1]
	}

	parts = strings.Split(url, "/")
	return parts[0]
}
