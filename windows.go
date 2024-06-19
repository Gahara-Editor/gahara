//go:build windows

package main

import "fmt"

func ExtractFFmpeg() (string, error) {
	return "", fmt.Errorf("unsupported platform")
}
