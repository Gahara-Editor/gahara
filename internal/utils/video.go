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
		return "", "", fmt.Errorf("Invalid file format, %s does not contain a file extension", fileName)
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

func FormatTime(seconds float64) string {
	duration := time.Duration(int64(seconds * float64(time.Second)))

	// Extract the whole seconds and remaining nanoseconds
	wholeSeconds := int(duration.Seconds())
	nanoseconds := duration.Nanoseconds() % int64(time.Second)

	// Format the time components as "HH:MM:SS"
	formattedTime := fmt.Sprintf("%02d:%02d:%02d", wholeSeconds/3600, (wholeSeconds%3600)/60, wholeSeconds%60)

	// Add milliseconds (rounded to six decimal places)
	formattedTime += fmt.Sprintf(".%06d", nanoseconds/1000)
	return formattedTime
}
