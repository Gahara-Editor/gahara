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
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/k1nho/gahara/ffmpegbuilder"
	"github.com/k1nho/gahara/internal/video"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Video struct {
	// ID: the unique identifier of the video
	ID string `json:"id"`
	// Name: the name of the video file (includes extension)
	Name string `json:"name"`
	// Extension: the container type of the video (mp4, avi, etc)
	Extension string `json:"extension"`
	// FilePath: the absolute path of the video
	FilePath string `json:"filepath"`
	// Duration: the duration of the video in seconds
	Duration float64 `json:"duration"`
}

type Interval struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type VideoProcessingResult struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MonitoringOpts struct {
	terms map[string]bool
}

func NewVideo(name string, extension string, filepath string, duration float64) *Video {
	return &Video{
		ID:        strings.Replace(uuid.New().String(), "-", "", -1),
		Name:      name,
		Extension: extension,
		FilePath:  filepath,
		Duration:  duration,
	}
}

func NewMonitoringOpts(observeParams ...string) *MonitoringOpts {
	terms := make(map[string]bool)
	for _, word := range observeParams {
		terms[word] = true
	}
	return &MonitoringOpts{
		terms: terms,
	}
}

func NewVideoProcessingResult(id string, name string, status string, msg string) *VideoProcessingResult {
	if id == "" {
		id = strings.Replace(uuid.New().String(), "-", "", -1)
	}

	return &VideoProcessingResult{
		ID:      id,
		Name:    name,
		Status:  status,
		Message: msg,
	}
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
		pfile := NewVideo(name, filepath.Ext(proxyFile), a.config.ProjectDir, 0)

		cancelDurationListener := wruntime.EventsOnce(a.ctx, video.EVT_DURATION_EXTRACTED, func(duration ...interface{}) {
			pfile.Duration = duration[0].(float64)
			wruntime.EventsEmit(a.ctx, video.EVT_PROXY_FILE_CREATED, pfile)
		})
		defer cancelDurationListener()

		err := a.FFmpegQuery(video.QUERY_CREATE_PROXY_FILE, video.ProcessingOpts{
			Filename:    name,
			VideoFormat: fmt.Sprintf(".%s", ext),
			InputPath:   filepath.Dir(inputFilePath),
			OutputPath:  a.config.ProjectDir,
		})
		if err != nil {
			wruntime.LogError(a.ctx, fmt.Sprintf("could not create the proxy file for %s: %s", inputFilePath, err.Error()))
			wruntime.EventsEmit(a.ctx, video.EVT_PROXY_ERROR_MSG, fmt.Sprintf("failed to import %s", fileName))
			return
		}

		wruntime.LogInfo(a.ctx, fmt.Sprintf("proxy file created: %s", fileName))
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

func (a *App) GetProjectThumbnail(projectName string) (string, error) {
	thumbnailDir := path.Join(a.config.GaharaDir, projectName)
	projectDir, err := os.Open(thumbnailDir)
	if err != nil {
		wruntime.LogError(a.ctx, "directory does not exists for the project")
		return "", err
	}
	defer projectDir.Close()

	files, err := projectDir.ReadDir(0)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the files of the project")
		return "", err
	}

	var thumbnailPath string
	for _, project := range files {
		if !project.IsDir() && filepath.Ext(project.Name()) == ".png" {
			thumbnailPath = path.Join(thumbnailDir, project.Name())
			break
		}
	}

	if thumbnailPath == "" {
		return thumbnailPath, fmt.Errorf("no thumbnail found")
	}
	return thumbnailPath, nil

}

func (a *App) SaveProjectFiles(projectFiles []Video) error {
	data, err := json.MarshalIndent(projectFiles, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(a.config.ProjectDir, "metadata.json"), data, 0644)
	if err != nil {
		return err
	}

	wruntime.LogInfo(a.ctx, "project files have been saved")
	return nil
}

// SaveTimeline: save project timeline into the project filesystem
func (a *App) SaveTimeline() error {
	if a.Timeline.VideoNodes == nil && len(a.Timeline.VideoNodes) <= 0 {
		return fmt.Errorf("timeline is empty, could not save timeline")
	}
	data, err := json.MarshalIndent(a.Timeline, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(a.config.ProjectDir, "timeline.json"), data, 0644)
	if err != nil {
		return err
	}
	wruntime.LogInfo(a.ctx, fmt.Sprintf("%s: timeline has been saved", time.Now().String()))
	return nil
}

// LoadTimeline: retrieve saved project timeline, if any, from filesystem
func (a *App) LoadTimeline() (video.Timeline, error) {
	var timeline video.Timeline
	timelinePath := path.Join(a.config.ProjectDir, "timeline.json")
	if _, err := os.Stat(timelinePath); err != nil {
		return timeline, fmt.Errorf("no timeline found for this project")
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

	if len(a.Timeline.VideoNodes) == 0 {
		wruntime.LogInfo(a.ctx, "empty timeline")
		return timeline, fmt.Errorf("empty timeline")
	}

	wruntime.LogInfo(a.ctx, "timeline has been loaded!")
	return a.GetTimeline(), nil
}

// LoadTimeline: retrieves saved project files, if any, from filesystem
func (a *App) LoadProjectFiles() ([]Video, error) {
	var videoFiles []Video
	metadataPath := path.Join(a.config.ProjectDir, "metadata.json")
	if _, err := os.Stat(metadataPath); err != nil {
		return videoFiles, fmt.Errorf("No video files found for this project")
	}

	bytes, err := os.ReadFile(metadataPath)
	if err != nil {
		wruntime.LogError(a.ctx, "could not read the video files metadata file")
		return videoFiles, fmt.Errorf("could not read timeline file")
	}

	err = json.Unmarshal(bytes, &videoFiles)
	if err != nil {
		wruntime.LogError(a.ctx, "could not unmarshal the files")
		return videoFiles, err
	}

	if len(videoFiles) == 0 {
		wruntime.LogInfo(a.ctx, "empty video files")
		return videoFiles, fmt.Errorf("empty video files")
	}

	wruntime.LogInfo(a.ctx, "video files loaded")
	return videoFiles, nil

}

// GetTimeline: returns the video timeline which is composed of video nodes
func (a *App) GetTimeline() video.Timeline {
	return a.Timeline
}

// InsertInterval: inserts a video node with some interval [a,b]
func (a *App) InsertInterval(rid string, name string, start, end float64, pos int) (video.VideoNode, error) {
	return a.Timeline.Insert(rid, name, start, end, pos)
}

// RemoveInterval: removes a video node with some interval [a,b]
func (a *App) RemoveInterval(pos int) error {
	return a.Timeline.Delete(pos)
}

// SplitInterval: splits a video node with some interval [a,b].
func (a *App) SplitInterval(eventType string, pos int, start, end float64) ([]video.VideoNode, error) {
	return a.Timeline.Split(eventType, pos, start, end)
}

// DeleteRIDReferences: removes all timeline references of a root id
func (a *App) DeleteRIDReferences(rid string) error {
	return a.Timeline.DeleteRIDReferences(rid)
}

func (a *App) RenameVideoNode(pos int, name string) error {
	return a.Timeline.RenameVideoNode(pos, name)
}

func (a *App) ToggleLossless(pos int) error {
	return a.Timeline.ToggleLossless(pos)
}

func (a *App) MarkAllLossless() error {
	return a.Timeline.MarkAllLossless()
}

func (a *App) UnmarkAllLossless() error {
	return a.Timeline.UnmarkAllLossless()
}

// ResetTimeline: cleanup timeline state in memory
func (a *App) ResetTimeline() {
	a.Timeline = video.NewTimeline()
}

// GetTrackDuration: retrieves the total video duration of a track
func (a *App) GetTrackDuration() (float64, error) {
	if a.Timeline.VideoNodes == nil {
		return 0, fmt.Errorf("no timeline exists")
	}

	duration := 0.0
	for _, videoNode := range a.Timeline.VideoNodes {
		duration += videoNode.End - videoNode.Start
	}
	return duration, nil
}

// GetOutputFileSavePath: retrieves the output path where the resulting video should be saved
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

// FFmpegQuery: produces an asset (video, image) with FFmpeg given a query type, and processing opts
func (a *App) FFmpegQuery(queryType string, userOpts video.ProcessingOpts) error {
	defer wruntime.EventsEmit(a.ctx, video.EVT_FFMPEG_EXEC_ENDED)

	if userOpts.InputPath == "" {
		userOpts.InputPath = a.config.ProjectDir
	}
	if err := userOpts.ValidateRequiredFields(queryType); err != nil {
		return err
	}

	switch queryType {
	case video.QUERY_FILTERGRAPH:
		if err := a.queryFiltergraph(userOpts); err != nil {
			return err
		}
	case video.QUERY_LOSSLESS_CUT:
		if err := a.queryLosslessCut(userOpts); err != nil {
			return err
		}
	case video.QUERY_CREATE_PROXY_FILE:
		if err := a.queryCreateProxyFile(userOpts); err != nil {
			return err
		}
	case video.QUERY_CREATE_THUMBNAIL:
		if err := a.queryCreateThumbnail(userOpts); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid query type (q_filtergraph, q_create_proxy_file, q_create_thumbnail, q_lossless_cut)")
	}
	return nil
}

// queryFiltergraph: executes a filtergraph query, currently merge clips
func (a *App) queryFiltergraph(userOpts video.ProcessingOpts) error {
	query, err := ffmpegbuilder.MergeClipsQuery(a.Timeline.VideoNodes, userOpts)
	if err != nil {
		return err
	}
	err = a.executeFFmpegQuery(query, NewMonitoringOpts(video.OBV_OUT_TIME_US))
	if err != nil {
		return err
	}

	wruntime.EventsEmit(a.ctx, video.EVT_FFMPEG_RESULT, NewVideoProcessingResult("", userOpts.Filename, Success, ffmpegbuilder.GetFullOutputPath(userOpts)))
	wruntime.EventsEmit(a.ctx, video.EVT_EXPORT_MSG, fmt.Sprintf("Finished exporting %s%s", userOpts.Filename, userOpts.VideoFormat))
	return nil
}

// queryLosslessCut: executes LosslessCut for a batch of video nodes
func (a *App) queryLosslessCut(userOpts video.ProcessingOpts) error {
	var (
		wg         = new(sync.WaitGroup)
		msgChannel = make(chan VideoProcessingResult)
	)

	go func() {
		defer close(msgChannel)
		for msg := range msgChannel {
			wruntime.EventsEmit(a.ctx, video.EVT_FFMPEG_RESULT, msg)
		}
	}()

	for _, videoNode := range a.Timeline.VideoNodes {
		if !videoNode.LosslessExport {
			continue
		}
		wg.Add(1)
		go func(vNode video.VideoNode) {
			defer wg.Done()
			userOpts.Filename = vNode.Name
			query, err := ffmpegbuilder.LosslessCutQuery(vNode, userOpts)
			if err != nil {
				msgChannel <- VideoProcessingResult{ID: vNode.ID, Status: Failed, Message: err.Error()}
				return
			}
			err = a.executeFFmpegQuery(query, nil)
			if err != nil {
				msgChannel <- VideoProcessingResult{ID: vNode.ID, Name: vNode.Name, Status: Failed, Message: err.Error()}
				return
			}
			msgChannel <- VideoProcessingResult{ID: vNode.ID, Name: vNode.Name, Status: Success, Message: ffmpegbuilder.GetFullOutputPath(userOpts)}
		}(videoNode)
	}
	wg.Wait()
	return nil
}

// queryCreateProxyFile: executes a conversion query for the given video
func (a *App) queryCreateProxyFile(userOpts video.ProcessingOpts) error {
	query, err := ffmpegbuilder.CreateProxyFileQuery(userOpts, ".mov")
	if err != nil {
		return err
	}

	err = a.executeFFmpegQuery(query, NewMonitoringOpts(video.OBV_OUT_TIME))
	if err != nil {
		return err
	}

	return nil
}

// queryCreateThumbnail: executes a query to generate a thumbnail for a video (picks 1st frame)
func (a *App) queryCreateThumbnail(userOpts video.ProcessingOpts) error {
	query, err := ffmpegbuilder.CreateThumbnailQuery(userOpts, ".png")
	if err != nil {
		return err
	}

	err = a.executeFFmpegQuery(query, nil)
	if err != nil {
		return err
	}

	return nil
}

// executeFFmpegQuery: executes an ffmpeg query
func (a *App) executeFFmpegQuery(query string, monitoringOpts *MonitoringOpts) error {
	// TODO: implement windows
	cmd := exec.Command("bash", "-c", query)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not initialize ffmpeg monitoring")
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("could not initialize video export")
	}
	if monitoringOpts != nil {
		go a.monitorFFmpegOuput(stderrPipe, monitoringOpts)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("could not export the video")
	}
	return nil
}

// monitorFFmpegOuput: monitors ffmpeg query progress
func (a *App) monitorFFmpegOuput(FFmpegOut io.ReadCloser, monitoringOpts *MonitoringOpts) {
	wruntime.LogInfo(a.ctx, "monitoring FFmpeg query")
	total, err := a.GetTrackDuration()
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(FFmpegOut)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, video.OBV_OUT_TIME_US) && monitoringOpts.terms[video.OBV_OUT_TIME_US] {
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
		if strings.Contains(line, video.OBV_OUT_TIME) && monitoringOpts.terms[video.OBV_OUT_TIME] {
			args := strings.Split(line, "=")
			if args[0] != video.OBV_OUT_TIME {
				continue
			}
			duration, err := convertHMStoSeconds(args[1])
			if err != nil {
				// TODO: handle conversion error
				continue
			}
			wruntime.EventsEmit(a.ctx, video.EVT_DURATION_EXTRACTED, duration)
		}
	}
}

func convertHMStoSeconds(hms string) (float64, error) {
	parts := strings.Split(hms, ".")
	if len(parts) != 2 {
		return 0.0, fmt.Errorf("could not parse decimal part")
	}

	timeParts := strings.Split(parts[0], ":")
	if len(timeParts) != 3 {
		return 0.0, fmt.Errorf("could not parse time part")
	}

	hours, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return 0.0, fmt.Errorf("could not convert hours to int")
	}

	minutes, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return 0.0, fmt.Errorf("could not convert minutes to int")
	}

	seconds, err := strconv.Atoi(timeParts[2])
	if err != nil {
		return 0.0, fmt.Errorf("could not convert seconds to int")
	}

	microseconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0.0, fmt.Errorf("could not convert microseconds to int")
	}

	totalSeconds := float64(hours*3600+minutes*60+seconds) + float64(microseconds)/1000000
	return totalSeconds, nil
}

func getVideoDuration(userOpts video.ProcessingOpts) (float64, error) {
	query, err := ffmpegbuilder.CheckVideoDuration(userOpts)
	if err != nil {
		return 0, err
	}
	cmd := exec.Command("bash", "-c", query)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 0, fmt.Errorf("could not initialize ffmpeg monitoring")
	}

	scanner := bufio.NewScanner(stderrPipe)

	err = cmd.Start()
	if err != nil {
		return 0, fmt.Errorf("could not initialize video duration extraction")
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, video.OBV_DURATION) {
			hms := strings.Split(strings.TrimSpace(strings.Split(line, ",")[0]), "Duration: ")[1]
			duration, err := convertHMStoSeconds(hms)
			if err != nil {
				continue
			}
			return duration, nil
		}

	}

	err = cmd.Wait()
	if err != nil {
		return 0, fmt.Errorf("could not extract duration of the video")
	}
	return 0, nil

}
