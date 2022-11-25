package cliHttp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func IsCommandFileURL(commandFilePath string) bool {
	_, err := url.ParseRequestURI(commandFilePath)
	return !(err == nil)
}

func DownloadFile(urlString string) (bool, error) {
	fileURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Cannot download file: ", urlString)
		return false, err
	}

	path := fileURL.Path
	urlSegments := strings.Split(path, "/")
	outputFileName := urlSegments[len(urlSegments)-1]
	outputFilePath := fmt.Sprintf("%s%s", "~/.climanager/", outputFileName)

	response, err := http.Get(path)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	fileHandle, err := os.OpenFile(outputFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return false, err
	}
	defer fileHandle.Close()

	_, err = io.Copy(fileHandle, response.Body)
	if err != nil {
		fmt.Println("Could not write to file: ", outputFilePath)
		return false, err
	}

	return true, nil
}
