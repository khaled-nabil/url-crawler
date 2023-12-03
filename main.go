package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
	for _, link := range urls {
		// Extract domain name from URL
		domain := urlToFilename(link)

		// Make a GET request
		resp, err := http.Get(link)
		if err != nil {
			fmt.Println("Error making HTTP request:", err, link)
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		// Save the response body to a file
		fileName := fmt.Sprintf("%s/%s.txt", outputFolderPath, domain)
		err = os.WriteFile(fileName, body, os.ModePerm)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		fmt.Printf("Saved %s to %s\n", link, fileName)
	}
}

func urlToFilename(inputURL string) string {
	// Parse the URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		// Handle error
		return ""
	}

	suffix := filepath.Base(parsedURL.Path)

	invalidCharRegex := regexp.MustCompile(`[^\w-]`)
	safeHostname := invalidCharRegex.ReplaceAllString(parsedURL.Hostname(), "_")
	safeSuffix := invalidCharRegex.ReplaceAllString(suffix, "_")

	filename := safeHostname + "_" + safeSuffix

	maxLength := 50
	if len(filename) > maxLength {
		filename = filename[:maxLength]
	}

	return filename
}
