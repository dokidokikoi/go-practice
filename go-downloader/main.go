package main

import (
	"fmt"
	"os"
	"strings"

	"practice/go-downloader/download"
)

func main() {
	fmt.Println("Download Started")

	fileUrl := os.Args[1]
	var filePath string

	if len(os.Args) < 3 {
		filePath = strings.Split(fileUrl, "/")[len(strings.Split(fileUrl, "/"))-1]
	} else {
		filePath = os.Args[2]
		if filePath == "" {
			filePath = strings.Split(fileUrl, "/")[len(strings.Split(fileUrl, "/"))-1]
		}
	}

	err := download.DownloadFile(filePath, fileUrl)
	if err != nil {
		panic(err)
	}

	fmt.Println("Download Finished")
}
