package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
func (a *App) createProxyFile(inputFilePath string) {
	if inputFilePath == "" {
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, "no file selected")
		return
	}

	fileName := video.GetFilename(inputFilePath)
	name, ext, err := video.GetNameAndExtension(fileName)
	if err != nil {
		wruntime.LogError(a.ctx, "invalid file format")
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, "invalid file format")
		return
	}

	if !video.IsValidExtension("." + ext) {
		wruntime.LogError(a.ctx, "invalid file extension")
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, "invalid file extension")
		return
	}

	proxyFile := fmt.Sprintf("%s.mov", name)
	pathProxyFile := path.Join(a.config.ProjectDir, proxyFile)

	// check that a proxy has not already been created for the file
	_, err = os.Stat(pathProxyFile)
	if os.IsNotExist(err) {
		cmd := video.CreateProxyFileCMD(inputFilePath, pathProxyFile)
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_PIPELINE_MSG, name)
		err := cmd.Run()
		if err != nil {
			wruntime.LogError(a.ctx, fmt.Sprintf("could not create the proxy file for %s: %s", inputFilePath, err.Error()))
			wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, fmt.Sprintf("failed to import %s", fileName))
			return
		}

		wruntime.LogInfo(a.ctx, fmt.Sprintf("proxy file created: %s", fileName))
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_FILE_CREATED,
			Video{Name: name, FilePath: a.config.ProjectDir, Extension: filepath.Ext(fileName)})

		// Once the proxy file is created, generate a thumbnail
		err = a.GenerateThumbnail(pathProxyFile)
		if err != nil {
			wruntime.LogError(a.ctx, fmt.Sprintf("could not generate thumbnail for proxy file %s: %s", inputFilePath, err.Error()))
		}
		return
	} else if err != nil {
		wruntime.LogError(a.ctx, fmt.Sprintf("file finding error: %s", err.Error()))
		wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, fmt.Sprintf("failed to import %s", fileName))
		return
	}

	wruntime.LogInfo(a.ctx, fmt.Sprintf("proxy file found: %s", fileName))
	wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, fmt.Sprintf("file is already in project %s", fileName))

}

// GenerateThumbnail: given an input file, generates a single frame that can be used as thumbnail
func (a *App) GenerateThumbnail(inputFilePath string) error {
	inputFile, err := os.Stat(inputFilePath)
	if err != nil {
		return fmt.Errorf("could not find proxy file: %s", err.Error())
	}

	filename, _, err := video.GetNameAndExtension(inputFile.Name())
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
	wruntime.LogInfo(a.ctx, fmt.Sprintf("%s: Timeline has been saved", time.Now().String()))
	return nil
}

// LoadTimeline: retrieve saved project timeline, if any, from filesystem
func (a *App) LoadTimeline() (video.Timeline, error) {
	var timeline video.Timeline
	timelinePath := path.Join(a.config.ProjectDir, "timeline.json")
	if _, err := os.Stat(timelinePath); err != nil {
		return timeline, fmt.Errorf("No timeline found for this project")
	}

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

// ResetTimeline: cleanup timeline state in memory
func (a *App) ResetTimeline() {
	a.Timeline = video.NewTimeline()
}

func (a *App) GetTrackDuration() (float64, error) {
	if a.Timeline.VideoNodes == nil || len(a.Timeline.VideoNodes) == 0 {
		return 0, fmt.Errorf("no timeline exists")
	}

	duration := 0.0
	for _, videoNode := range a.Timeline.VideoNodes {
		duration += videoNode.End - videoNode.Start
	}
	return duration, nil
}

// FFmpegQueryBuild: builds the FFmpeg query based on timeline state and processing options (codec, resolution)
func (a *App) ffmpegQueryBuild(processingOpts *video.ProcessingOpts) (string, error) {
	err := video.ValidateProcessingOpts(processingOpts)
	if err != nil {
		return "", err
	}

	inputArgs, err := a.Timeline.InputArgs()
	if err != nil {
		return "", err
	}

	query, err := a.Timeline.MergeClipsQuery(processingOpts)
	if err != nil {
		return "", err
	}

	outputArgs, err := a.Timeline.OutputArgs(processingOpts)
	if err != nil {
		return "", err
	}

	var ffmpegQuery strings.Builder
	ffmpegQuery.WriteString("ffmpeg -v quiet -stats_period 5s -progress pipe:2")
	ffmpegQuery.WriteString(inputArgs)
	ffmpegQuery.WriteString(query)
	ffmpegQuery.WriteString(outputArgs)
	return ffmpegQuery.String(), nil
}

func (a *App) GetOutputFileSavePath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get the user home directory")
	}

	saveFilepath, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		DefaultDirectory:           path.Join(hd),
		Title:                      "Export video",
		ShowHiddenFiles:            false,
		CanCreateDirectories:       true,
		TreatPackagesAsDirectories: true,
	})
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return saveFilepath, nil
}

func (a *App) ExportVideo(userOpts video.ProcessingOpts) (string, error) {
	query, err := a.ffmpegQueryBuild(&userOpts)
	if err != nil {
		return "", err
	}
	// TODO: handle win implementation
	cmd := exec.Command("bash", "-c", query)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("could not initialize ffmpeg monitoring")
	}

	err = cmd.Start()
	if err != nil {
		return "", fmt.Errorf("could not initialize video export")
	}
	go a.monitorFFmpegOuput(stderrPipe)

	err = cmd.Wait()
	if err != nil {
		return "", fmt.Errorf("could not export the video: %s%s", userOpts.Filename, userOpts.VideoFormat)
	}
	return fmt.Sprintf("%s%s has been processed!", userOpts.Filename, userOpts.VideoFormat), err
}

func (a *App) monitorFFmpegOuput(FFmpegOut io.ReadCloser) {
	total, err := a.GetTrackDuration()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(FFmpegOut)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "out_time_us") {
			args := strings.Split(line, "=")
			timeMicro, err := strconv.Atoi(args[1])
			if err != nil {
				continue
			}
			timeSeconds := (timeMicro / 1000000)
			if timeSeconds < 0 {
				continue
			}
			wruntime.EventsEmit(a.ctx, video.EVT_ENCODING_PROGRESS, (timeSeconds*100)/int(total))
		}
	}
}
