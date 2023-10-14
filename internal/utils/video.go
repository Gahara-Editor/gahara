package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/k1nho/gahara/internal/constants"
)

func GetFilename(path string) string {
	pathSlice := strings.Split(path, "/")
	return pathSlice[len(pathSlice)-1]
}

func GetNameAndExtension(fileName string) (name string, ext string, err error) {
	fileSlice := strings.Split(fileName, ".")
	if len(fileSlice) != 2 {
		return "", "", fmt.Errorf("Invalid file format does not contain a file extension")
	}
	return fileSlice[0], fileSlice[1], nil
}

func IsValidExtension(extension string) bool {
	for _, ext := range constants.ValidVideoExtensions {
		if extension == ext {
			return true
		}
	}
	return false
}

func FormatTime(t time.Duration) string {
	duration := time.Duration(t) * time.Second
	// Extract hours, minutes, and seconds
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) - hours*60
	seconds := int(duration.Seconds()) - hours*3600 - minutes*60

	// Format the duration as "hr:min:sec"
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
