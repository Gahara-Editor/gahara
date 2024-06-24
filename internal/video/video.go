package video

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getValidVideoExtensions() []string {
	return []string{".mov", ".mkv", ".mp4", ".m4v", ".avi", ".wmv", ".flv", ".h264", ".hevc", ".3gp", ".3g2", ".ogv", ".264", ".265", ".webm"}
}

func getCompatibleOGV() []string {
	return []string{CODEC_THEORA, CODEC_COPY}
}
func getCompatibleWebm() []string {
	return []string{CODEC_VP9, CODEC_COPY}
}

func getCompatibleRest() []string {
	return []string{CODEC_H264, CODEC_H264_RGB, CODEC_H265, CODEC_COPY}
}

const (
	// high order query types
	QUERY_FILTERGRAPH       = "q_filtergraph"
	QUERY_LOSSLESS_CUT      = "q_lossless_cut"
	QUERY_CREATE_PROXY_FILE = "q_create_proxy_file"
	QUERY_CREATE_THUMBNAIL  = "q_create_thumbnail"
	// Epsilon: margin for floating point checks
	Epsilon          = 1e-6
	EVT_CHANGE_ROUTE = "evt_change_route"
	// EVT_FFMPEG_RESULT: signals the result of a FFmpeg query in the form of a VideoProcessingResult
	EVT_FFMPEG_RESULT = "evt_ffmpeg_result"
	// EVT_FFMPEG_EXEC_ENDED: signals that the ffmpeg query has ended (does not indicate success of the query)
	EVT_FFMPEG_EXEC_ENDED = "evt_ffmpeg_exec_ended"
	// EVT_ENABLE_VIM_MODE
	EVT_TOGGLE_VIM_MODE = "evt_toggle_vim_mode"
	// EVT_CHANGE_VIM_MODE: sets the vim mode in timeline (select, remove, timeline)
	EVT_CHANGE_VIM_MODE = "evt_change_vim_mode"
	// EVT_SAVED_TIMELINE: triggers once the timeline has been saved
	EVT_SAVED_TIMELINE = "evt_saved_timeline"
	//EVT_ZOOM_TIMELINE: triggers a zoom for the timeline
	EVT_ZOOM_TIMELINE = "evt_zoom_timeline"
	// EVT_SPLITCLIP_EDIT: invokes the handleTwoCut function
	EVT_SPLITCLIP_EDIT = "evt_splitclip_edit"
	// EVT_INSERTCLIP_EDIT: triggers clip insertion (yank->paste or search list)
	EVT_INSERTCLIP_EDIT = "evt_insertclip_edit"
	//EVT_TOGGLE_LOSSLESS: toggles the video node status to be exported as lossless
	EVT_TOGGLE_LOSSLESS = "evt_toggle_lossless"
	// EVT_MARK_ALL_LOSSLESS: marks all the video nodes as lossless
	EVT_MARK_ALL_LOSSLESS = "evt_mark_all_lossless"
	//EVT_UNMARK_ALL_LOSSLESS: umarks all the video nodes marked as lossless
	EVT_UNMARK_ALL_LOSSLESS = "evt_unmark_all_lossless"
	// EVT_EXECUTE_EDIT: execute the current edit
	EVT_EXECUTE_EDIT = "evt_execute_edit"
	// EVT_PLAY_TRACK: plays the clips on the track (starting from current pos and clip time)
	EVT_PLAY_TRACK = "evt_play_track"
	//EVT_TRACK_MOVE: move vim cursor on the current track
	EVT_TRACK_MOVE = "evt_track_move"
	//EVT_UPLOAD_FILE: triggers the native upload file
	EVT_UPLOAD_FILE = "evt_upload_file"
	//EVT_OPEN_SEARCH_LIST: opens a search list containing all the uploaded files
	EVT_OPEN_SEARCH_LIST = "evt_open_search_list"
	//EVT_YANK_CLIP: copies the selected video node on track
	EVT_YANK_CLIP = "evt_yank_clip"
	// EVT_OPEN_RENAME_CLIP_MODAL: opens the rename clip modal
	EVT_OPEN_RENAME_CLIP_MODAL = "evt_open_rename_clip_modal"
	// EVT_INTERVAL_CUT: interval cut event
	EVT_INTERVAL_CUT = "intervalCut"
	// EVT_SLICE_CUT: slice cut event
	EVT_SLICE_CUT = "evt_slice_cut"
	// EVT_DURATION_EXTRACTED: duration extracted from ffmpeg execution
	EVT_DURATION_EXTRACTED = "evt_duration_extracted"
	// EVT_ENCODING_PROGRESS: encoding progress event
	EVT_ENCODING_PROGRESS = "evt_encoding_progress"
	// EVT_EXPORT_MSG: message for exporting a video events
	EVT_EXPORT_MSG = "evt_export_msg"
	// EVT_PIPELINE_MSG: proxy pipeline message event
	EVT_PROXY_PIPELINE_MSG = "evt_proxy_pipeline_msg"
	//EVT_PROXY_ERROR_MSG: error ocurred while creating proxy file
	EVT_PROXY_ERROR_MSG = "evt_proxy_error_msg"
	//EVT_ENCODING_PROGRESS: proxy file has been created event
	EVT_PROXY_FILE_CREATED = "evt_proxy_file_created"
	//SCALE_256x256: Resolution 256x256
	SCALE_256x256 = "256x256"
	//SCALE_316_192: Resolution 316x192
	SCALE_316_192 = "316x192"
	//SCALE_640x480: Resolution 640x480 (SD)
	SCALE_640x480 = "640x480"
	// SCALE_1280X720: Resolution 1280x720 (HD)
	SCALE_1280X720 = "1280x720"
	// SCALE_1920X1080: Resolution 1920x1080 (full HD)
	SCALE_1920X1080 = "1920x1080"
	// SCALE_1920X1080: Resolution 2560x1440 (QHD)
	SCALE_2560x1440 = "2560x1440"
	// SCALE_3840X2160: Resolution 3840x2160 (UHD)
	SCALE_3840X2160 = "3840x2160"
	// CODEC_H265: codec H.265
	CODEC_H265 = "libx265"
	// CODEC_H264: codec H.264
	CODEC_H264 = "libx264"
	// CODEC_H264_RGB: codec H.264 RGB
	CODEC_H264_RGB = "libx264rgb"
	// CODEC_VP9: codec VP9 RGB (.webm)
	CODEC_VP9 = "libvpx-vp9"
	// CODEC_VP9: codec theora (.ogv)
	CODEC_THEORA = "libtheora"
	// CODEC_COPY: codec copy for bitstream copy (skips re-encoding the video)
	CODEC_COPY = "copy"
	// CRF_23: constant rate factor 23, lower number usually means better quality
	CRF_23 = "23"
	// CRF_22: constant rate factor 22, lower number usually means better quality
	CRF_22 = "22"
	// CRF_21: constant rate factor 22, lower number usually means better quality
	CRF_21 = "21"
	// CRF_20: constant rate factor 22, lower number usually means better quality
	CRF_20 = "20"
	// CRF_19: constant rate factor 22, lower number usually means better quality
	CRF_19 = "19"
	// CRF_18: constant rate factor 18, lower number usually means better quality
	CRF_18 = "18"
	// PRESET_SLOW: slow preset provides better compression (quality per file size)
	PRESET_SLOW = "slow"
	// PRESET_MEDIUM: medium preset is a balance of encoding speed to compression
	PRESET_MEDIUM = "medium"
	// PRESET_FAST: fast preset provides faster encoding speed at the tradeoff of better compression
	PRESET_FAST = "fast"
	// OUT_TIME_US: out_time_us term to monitor in ffmpeg execution
	OBV_OUT_TIME_US = "out_time_us"
	// OUT_TIME: out_time term to monitor in ffmpeg execution
	OBV_OUT_TIME = "out_time"
	// Duration: duration term to monitor in ffmpeg execution
	OBV_DURATION = "Duration"
)

type VideoNode struct {
	// Start: the start of the interval
	Start float64 `json:"start"`
	// End: the end of the interval
	End float64 `json:"end"`
	// RID: the root ID of the node, that is, the original video from which this nodes derives
	RID string `json:"rid"`
	// ID: the ID of the video node
	ID string `json:"id"`
	// Name: the name given by the user to the clip
	Name string `json:"name"`
	// Lossless
	LosslessExport bool `json:"losslessexport"`
}

type Timeline struct {
	// VideoNodes: all the video nodes of the timeline
	VideoNodes []VideoNode `json:"video_nodes"`
}

type ThumbnailOpts struct {
	// StartTime: the second to get the thumbnail
	StartTime string `json:"start_time"`
	// Resolution: the size of the thumbnail
	Resolution string `json:"resolution"`
}

type ProcessingOpts struct {
	// Resolution (640x480,1280x720, 1920x1080, 2560x1440, 3840x2160)
	Resolution string `json:"resolution"`
	// Codec: video codec (libx264, libx265)
	Codec string `json:"codec"`
	// CRF: Constant Rate Factor
	CRF string `json:"crf"`
	// Preset: Encoding speed to compression ratio
	Preset string `json:"preset"`
	// InputPath: The file path of the input video
	InputPath string `json:"input_path,omitempty"`
	// OutputPath: The file path of the output video
	OutputPath string `json:"output_path"`
	// Filename: the name of the file
	Filename string `json:"filename"`
	// VideoFormat: the video format (.mov, .mp4)
	VideoFormat string `json:"video_format"`
}

func NewTimeline() Timeline {
	return Timeline{VideoNodes: []VideoNode{}}
}

func createVideoNode(rid string, name string, start, end float64) VideoNode {
	if name == "" {
		name = "Node"
	}
	return VideoNode{
		RID:   rid,
		ID:    strings.Replace(uuid.New().String(), "-", "", -1),
		Name:  name,
		Start: start,
		End:   end,
	}
}

func (tl *Timeline) Insert(rid string, name string, start, end float64, pos int) (VideoNode, error) {
	var videoNode VideoNode
	if pos < 0 || pos > len(tl.VideoNodes) {
		return videoNode, fmt.Errorf("insertion position %d is invalid", pos)
	}

	videoNode = createVideoNode(rid, name, start, end)
	tl.VideoNodes = slices.Insert(tl.VideoNodes, pos, videoNode)

	return videoNode, nil
}

func (tl *Timeline) Delete(pos int) error {
	if pos < 0 || pos >= len(tl.VideoNodes) {
		return fmt.Errorf("delete position is invalid")
	}
	if len(tl.VideoNodes) == 0 {
		return fmt.Errorf("there are no video clips to delete in track")
	}
	tl.VideoNodes = slices.Delete(tl.VideoNodes, pos, pos+1)
	return nil
}

func (tl *Timeline) RenameVideoNode(pos int, name string) error {
	if pos < 0 || pos >= len(tl.VideoNodes) {
		return fmt.Errorf("clip position invalid")
	}
	if len(tl.VideoNodes) == 0 {
		return fmt.Errorf("there are no video clips to rename in track")
	}
	tl.VideoNodes[pos].Name = name
	return nil
}

func (tl *Timeline) ToggleLossless(pos int) error {
	if pos < 0 || pos >= len(tl.VideoNodes) {
		return fmt.Errorf("clip position invalid %d", pos)
	}
	if len(tl.VideoNodes) == 0 {
		return fmt.Errorf("there are no video clips to mark")
	}
	tl.VideoNodes[pos].LosslessExport = !tl.VideoNodes[pos].LosslessExport
	return nil
}

func (tl *Timeline) MarkAllLossless() error {
	if len(tl.VideoNodes) == 0 {
		return fmt.Errorf("there are no video clips to mark")
	}

	for i := range tl.VideoNodes {
		tl.VideoNodes[i].LosslessExport = true
	}

	return nil
}

func (tl *Timeline) UnmarkAllLossless() error {
	if len(tl.VideoNodes) == 0 {
		return fmt.Errorf("there are no video clips to mark")
	}

	for i := range tl.VideoNodes {
		tl.VideoNodes[i].LosslessExport = false
	}

	return nil
}

func (tl *Timeline) Split(eventType string, pos int, start, end float64) ([]VideoNode, error) {
	nodes := []VideoNode{}
	if pos < 0 || pos >= len(tl.VideoNodes) {
		return nodes, fmt.Errorf("split position is invalid")
	}
	if len(tl.VideoNodes) == 0 {
		return nodes, fmt.Errorf("there are no video clips to split in track")
	}

	splitNode := tl.VideoNodes[pos]

	switch eventType {
	case EVT_SLICE_CUT:
		if end > splitNode.Start && end+0.1 < splitNode.End {
			nodes = append(nodes, createVideoNode(splitNode.RID, splitNode.Name, start, end), createVideoNode(splitNode.RID, splitNode.Name, end+0.1, splitNode.End))
		}
	case EVT_INTERVAL_CUT:
		if start-0.1 > splitNode.Start && end+0.1 < splitNode.End {
			nodes = append(nodes, createVideoNode(splitNode.RID, splitNode.Name, splitNode.Start, start-0.1), createVideoNode(splitNode.RID, splitNode.Name, start, end),
				createVideoNode(splitNode.RID, splitNode.Name, end+0.1, splitNode.End))
		}
	}

	if len(nodes) <= 0 {
		return nodes, fmt.Errorf("invalid cut range")
	}
	tl.VideoNodes = append(tl.VideoNodes[:pos], append(nodes, tl.VideoNodes[pos+1:]...)...)
	return nodes, nil
}

func (tl *Timeline) DeleteRIDReferences(rid string) error {
	if tl.VideoNodes == nil {
		return fmt.Errorf("no timeline exists")
	}
	tl.VideoNodes = slices.DeleteFunc(tl.VideoNodes, func(vn VideoNode) bool {
		return vn.RID == rid
	})
	return nil
}

// GenerateEditThumbnail: generate a thumbnail from a video
func GenerateEditThumb(inputFilePath string, outputFilePath string, opts ThumbnailOpts) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", inputFilePath, // input file
		"-vframes", "1", // pick 1 frame from the video
		"-s", SCALE_316_192, // scale of the video frame
		outputFilePath, // output file
	)
}

/*
GetFilename: return the file name from a path.
Ex: dir1/dir2/filename.txt -> filename.txt
*/
func GetFilename(path string) string {
	pathSlice := strings.Split(path, "/")
	return pathSlice[len(pathSlice)-1]
}

/*
GetNameAndExtension: return the name and extension from a filename.
Ex: filename.txt -> [filename, txt]
*/
func GetNameAndExtension(fileName string) (name string, ext string, err error) {
	fileSlice := strings.Split(fileName, ".")
	if len(fileSlice) != 2 {
		return "", "", fmt.Errorf("invalid file format, %s does not contain a file extension", fileName)
	}
	return fileSlice[0], fileSlice[1], nil
}

// IsValidExtension: checks if a given file extension is a video extension supported (.mov, .mp4)
func IsValidExtension(extension string) bool {
	for _, ext := range getValidVideoExtensions() {
		if extension == ext {
			return true
		}
	}
	return false
}

// FormatTime: given time in seconds, it returns time in format HH:MM:SS
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

func (p *ProcessingOpts) ValidateRequiredFields(queryType string) error {
	if p.Filename == "" {
		return fmt.Errorf("filename was not provided")
	}
	if p.VideoFormat == "" {
		return fmt.Errorf("video format was not provided")
	}
	for _, extension := range getValidVideoExtensions() {
		if strings.Contains(p.Filename, extension) {
			return fmt.Errorf(fmt.Sprintf("invalid filename %s", p.Filename))
		}
	}

	switch queryType {
	case QUERY_FILTERGRAPH:
		if p.OutputPath == "" {
			return fmt.Errorf("output path was not provided")
		}
		if !p.isCodecCompatible() {
			return fmt.Errorf("codec is not compatible with %s format", p.VideoFormat)
		}
	case QUERY_LOSSLESS_CUT:
		if p.OutputPath == "" {
			return fmt.Errorf("output path was not provided")
		}

	case QUERY_CREATE_PROXY_FILE, QUERY_CREATE_THUMBNAIL:
		if p.OutputPath == "" {
			return fmt.Errorf("output path was not provided")
		}
		if p.InputPath == "" {
			return fmt.Errorf("input path was not provided")
		}
	}
	return nil
}

func (p ProcessingOpts) isCodecCompatible() bool {
	switch p.VideoFormat {
	case ".webm":
		return slices.Contains(getCompatibleWebm(), p.Codec)
	case ".ogv":
		return slices.Contains(getCompatibleOGV(), p.Codec)
	default:
		return slices.Contains(getCompatibleRest(), p.Codec)
	}
}
