package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"

	"github.com/k1nho/gahara/internal/utils"
	"github.com/k1nho/gahara/internal/video"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Video struct {
	// Name: the name of the video file (includes extension)
	Name string `json:"name"`
	// Extension: the container type of the video (mp4, avi, etc)
	Extension string `json:"extension"`
	// FilePath: the absolute path of the video
	FilePath string `json:"filepath"`
}

type Interval struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
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
		err = a.GenerateThumbnail(pathProxyFile)
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
	startCut := utils.FormatTime(interval.Start)
	endCut := utils.FormatTime(interval.End)

	name, _, err := utils.GetNameAndExtension(path.Base(inputFile))
	if err != nil {
		return err
	}
	id := strings.Replace(uuid.New().String(), "-", "", -1)
	outputFile := path.Join(a.config.ProjectDir, fmt.Sprintf("%s_%s.mov", name, id))

	cmd := video.CutVideoInterval(inputFile, outputFile, startCut, endCut)
	err = cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("could not trim the video interval from %s to %s: %v", startCut, endCut, err)
		wruntime.LogError(a.ctx, errMsg)
		return fmt.Errorf(errMsg)
	}

	wruntime.LogInfo(a.ctx, fmt.Sprintf("video %s: [%s - %s] processed", name, startCut, endCut))
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
func (a *App) GenerateThumbnail(inputFilePath string) error {
	inputFile, err := os.Stat(inputFilePath)
	if err != nil {
		return fmt.Errorf("could not find proxy file: %s", err.Error())
	}

	filename, _, err := utils.GetNameAndExtension(inputFile.Name())
	if err != nil {
		return fmt.Errorf("could not extract name and extension")
	}
	// check that the thumbnail exists
	thumbnailPath := fmt.Sprintf("%s/%s.png", a.config.ProjectDir, filename)
	_, err = os.Stat(thumbnailPath)
	if err == nil {
		wruntime.LogInfo(a.ctx, fmt.Sprintf("thumbnail of video %s already exists", filename))
		return nil
	}

	cmd := video.GenerateEditThumb(inputFilePath, thumbnailPath, video.ThumbnailOpts{})
	err = cmd.Run()
	if err != nil {
		errMsg := fmt.Sprintf("could not generate the thumbnail for file %s: %s", filename, err.Error())
		wruntime.LogError(a.ctx, errMsg)
		return fmt.Errorf(errMsg)
	}

	wruntime.LogInfo(a.ctx, fmt.Sprintf("thumbnail for video: %s has been created", filename))
	return nil
}

// GetThumbnail: retrieve thumbnail for a given input file
func (a *App) GetThumbnail(inputFilePath string) error {
	filename := filepath.Base(inputFilePath)
	if filename == "." {
		return fmt.Errorf("file %s does not exists", filename)
	}

	filename = filename + ".png"
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not find thumbnail %s", filename)
	}
	defer file.Close()

	return nil

}

// SaveTimeline: save project timeline into the project filesystem
func (a *App) SaveTimeline() error {
	if a.Timeline.VideoNodes == nil && len(a.Timeline.VideoNodes) <= 0 {
		return fmt.Errorf("Timeline is empty, could not save timeline")
	}
	data, err := json.MarshalIndent(a.Timeline, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(a.config.ProjectDir, "timeline.json"), data, 0644)
	if err != nil {
		return err
	}
	wruntime.LogInfo(a.ctx, fmt.Sprintf("%+v", a.Timeline))
	return nil
}

// LoadTimeline: retrieve saved project timeline, if any, from filesystem
func (a *App) LoadTimeline() (video.Timeline, error) {
	var timeline video.Timeline
	timelinePath := path.Join(a.config.ProjectDir, "timeline.json")
	if _, err := os.Stat(timelinePath); err != nil {
		return timeline, fmt.Errorf("No timeline found for this project")
	}
	// the file exists, read it into the struct
	bytes, err := os.ReadFile(timelinePath)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the timeline file")
		return timeline, fmt.Errorf("could not read timeline file")
	}

	err = json.Unmarshal(bytes, &a.Timeline)
	if err != nil {
		wruntime.LogError(a.ctx, "could not unmarshal the timeline")
		return timeline, err
	}

	wruntime.LogInfo(a.ctx, "timeline file has been found!")
	return a.GetTimeline(), nil
}

// GetTimeline: returns the video timeline which is composed of video nodes
func (a *App) GetTimeline() video.Timeline {
	return a.Timeline
}

// InsertInterval: inserts a video node with some interval [a,b]
func (a *App) InsertInterval(rid string, start, end float64, pos int) (video.VideoNode, error) {
	return a.Timeline.Insert(rid, start, end, pos)
}

// RemoveInterval: removes a video node with some interval [a,b]
func (a *App) RemoveInterval(pos int) error {
	return a.Timeline.Delete(pos)
}

// SplitInterval: splits a video node with some interval [a,b].
func (a *App) SplitInterval(eventType string, pos int, start, end float64) ([]video.VideoNode, error) {
	return a.Timeline.Split(eventType, pos, start, end)
}

func (a *App) ResetTimeline() {
	a.Timeline = video.NewTimeline()
}
