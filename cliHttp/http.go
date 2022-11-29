package cliHttp

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func IsCommandFileURL(commandFilePath string) bool {
	url, err := url.ParseRequestURI(commandFilePath)
	if err != nil {
		return false
	}

	if len(url.Host) > 0 {
		return true
	}

	return false
}

func DownloadFile(urlString string) (*string, error) {
	fileURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Cannot download file: ", urlString)
		return nil, err
	}

	path := fileURL.Path
	urlSegments := strings.Split(path, "/")
	outputFileName := urlSegments[len(urlSegments)-1]

	userHomeDirectory, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home directory issue: ", err)
		return nil, err
	}

	outputFilePathDirectory := fmt.Sprintf("%s%s", userHomeDirectory, "/.climanager/")
	outputFilePath := fmt.Sprintf("%s%s", outputFilePathDirectory, outputFileName)

	if _, err := os.Stat(outputFilePathDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputFilePathDirectory, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	response, err := http.Get(urlString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	fileHandle, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return nil, err
	}
	defer fileHandle.Close()

	_, err = io.Copy(fileHandle, response.Body)
	if err != nil {
		fmt.Println("Could not write to file: ", outputFilePath)
		return nil, err
	}

	return &outputFilePath, nil
}
