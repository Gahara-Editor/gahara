//go:build darwin

package main

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed resources/darwin/ffmpeg
var FFmpegBinary embed.FS

func ExtractFFmpeg() (string, error) {
	data, err := FFmpegBinary.ReadFile("resources/darwin/ffmpeg")
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "ffmpeg")
	if err != nil {
		return "", err
	}

	ffmpegPath := filepath.Join(tempDir, "ffmpeg")
	err = os.WriteFile(ffmpegPath, data, 0755)
	if err != nil {
		return "", err
	}
	return ffmpegPath, nil
}
