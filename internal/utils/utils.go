package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func OpenLogFile(log_file_path string) (*os.File, error) {
	// Get log directory
	log_dir := log_file_path[:strings.LastIndex(log_file_path, "/")]
	// Create log directory, if isn't exists
	if _, err := os.Stat(log_dir); err != nil {
		os.MkdirAll(log_dir, 0766)
	}

	logFile, err := os.OpenFile(log_file_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

// Returns file bytes
func DownloadFileFromURL(url string) ([]byte, error) {
	// Download file
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while downloading %s - %w", url, err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	return bytes, err
}
func IsUserRoot() bool {
	usr_id := os.Getuid()
	if usr_id == -1 {
		fmt.Println("This program can only run in unix-like operating systems like linux and other...")
		os.Exit(1)
		return false
	} else if usr_id == 0 {
		return true
	} else {
		return false
	}
}
