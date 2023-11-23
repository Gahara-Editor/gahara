package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/k1nho/gahara/internal/utils"
	"github.com/k1nho/gahara/internal/video"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Interval struct {
	Start time.Duration `json:"start"`
	End   time.Duration `json:"end"`
}

type Video struct {
	Name     string `json:"name"`
	FilePath string `json:"filepath"`
}

// createProxyFile: creates the proxy file to be used for editing (preserve original media)
func (a *App) createProxyFile(inputFile, fileName string) {
	proxyFile := fmt.Sprintf("%s.mov", fileName)
	pathProxyFile := path.Join(a.config.ProjectDir, proxyFile)

	// check that a proxy has not already been created for the file
	_, err := os.Stat(pathProxyFile)
	if os.IsNotExist(err) {
		cmd := video.CreateProxyFileCMD(inputFile, pathProxyFile)
		err := cmd.Run()
		if err != nil {
			wruntime.LogError(a.ctx, fmt.Sprintf("could not create the proxy file for %s: %s", inputFile, err.Error()))
			return
		}

		// Once the proxy file is created, generate a thumbnail
		err = a.GenerateThumbnail(proxyFile)
		if err != nil {
			wruntime.LogError(a.ctx, fmt.Sprintf("could not generate thumbnail for proxy file %s: %s", inputFile, err.Error()))
		}
		wruntime.LogInfo(a.ctx, fmt.Sprintf("proxy file created: %s", fileName))
		return
	} else if err != nil {
		wruntime.LogError(a.ctx, fmt.Sprintf("file finding error: %s", err.Error()))
	}

	wruntime.LogInfo(a.ctx, fmt.Sprintf("proxy file found: %s", fileName))

}

// TrimVideoInterval: given an input file, and an interval (start,end), it returns the video with interval (start,end) removed
func (a *App) TrimVideoInterval(inputFile string, interval Interval) error {
	cmd := video.CutVideoInterval(inputFile, interval.Start, interval.End)
	err := cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("could not trim the video interval from %d to %d: %s", interval.Start, interval.End, err.Error())
		wruntime.LogError(a.ctx, errMsg)
		return fmt.Errorf(errMsg)
	}
	return nil
}

// GenerateVideoConcatFile: generates a .txt file with all the names of the video files to concatenate
func (a *App) GenerateVideoConcatFile(filenames []string) error {

	id, err := uuid.NewRandom()
	if err != nil {
		wruntime.LogError(a.ctx, "could not generate uuid for concat file")
		return err
	}

	concatFilePath := path.Join(a.config.ProjectDir, id.String()+".txt")
	concatFile, err := os.Create(concatFilePath)
	if err != nil {
		wruntime.LogError(a.ctx, "could not generate file.txt concatenation")
		return err
	}
	defer concatFile.Close()

	content := strings.Join(filenames, "\n")
	err = os.WriteFile(concatFilePath, []byte(content), 0644)
	if err != nil {
		wruntime.LogError(a.ctx, "could not write the file.txt concatenation")
		return err
	}

	return nil
}

// GenerateThumbnail: given an input file, generates a single frame that can be used as thumbnail
func (a *App) GenerateThumbnail(inputFile string) error {
	filename, _, err := utils.GetNameAndExtension(inputFile)
	if err != nil {
		return fmt.Errorf("could not find proxy file %s: %s", inputFile, err.Error())
	}

	// check that the thumbnail exists
	_, err = os.Stat(filename + ".png")
	if err == nil {
		return nil
	}

	cmd := video.GenerateEditThumb(inputFile, video.ThumbnailOpts{})
	err = cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("could not generate the thumbnail for file %s: %s", inputFile, err.Error())
		wruntime.LogError(a.ctx, errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}

func (a *App) GetThumbnail(inputFile string) error {
	filename := filepath.Base(inputFile)
	if filename == "." {
		return fmt.Errorf("file %s does not exists", inputFile)
	}

	filename = filename + ".png"
	_, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not find thumbnail %s", filename)
	}

	return nil

}
