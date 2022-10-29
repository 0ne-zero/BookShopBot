package utils

import (
	"crypto/rand"
	"encoding/hex"
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

func WriteBytesToFile(path string, bytes []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)
	return err
}

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
func IsFloatNumberRound(n float64) bool {
	str_n := fmt.Sprint(n)
	if strings.Contains(str_n, ".") {
		return false
	} else {
		return true
	}
}
func GenerateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return bytes, err
}
func GenerateRandomHex(length int) (string, error) {
	var byte_size = length
	if length%2 != 0 {
		byte_size += 1
	}
	bytes, err := GenerateRandomBytes(byte_size / 2)
	if err != nil {
		return "", err
	}
	hex := hex.EncodeToString(bytes)
	hex_len := len(hex)
	for hex_len != length {
		hex = hex[:hex_len-1]
		hex_len = len(hex)
	}
	return hex, nil
}
