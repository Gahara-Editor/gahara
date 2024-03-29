package video

import (
	"fmt"
	"os/exec"
	"path"
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

	// Epsilon: margin for floating point checks
	Epsilon = 1e-6
	// EVT_INTERVAL_CUT: interval cut event
	EVT_INTERVAL_CUT = "intervalCut"
	// EVT_SLICE_CUT: slice cut event
	EVT_SLICE_CUT = "evt_slice_cut"
	// EVT_ENCODING_PROGRESS: encoding progress event
	EVT_ENCODING_PROGRESS = "evt_encoding_progress"
	// EVT_PIPELINE_MSG: proxy pipeline message event
	EVT_PROXY_PIPELINE_MSG = "evt_proxy_pipeline_msg"
	//EVT_PROXY_ERROR_MSG: error ocurred while creating proxy file
	EVT_PROXY_ERROR_MSG = "evt_proxy_error_msg"
	//EVT_ENCODING_PROGRESS: proxy file has been created event
	EVT_PROXY_FILE_CREATED = "evt_proxy_file_created"
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
)

type VideoNode struct {
	// RID: the root ID of the node, that is, the original video from which this nodes derives
	RID string `json:"rid"`
	// ID: the ID of the video node
	ID string `json:"id"`
	// Start: the start of the interval
	Start float64 `json:"start"`
	// End: the end of the interval
	End float64 `json:"end"`
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

func createVideoNode(rid string, start, end float64) VideoNode {
	return VideoNode{
		RID:   rid,
		ID:    strings.Replace(uuid.New().String(), "-", "", -1),
		Start: start,
		End:   end,
	}
}

func (tl *Timeline) Insert(rid string, start, end float64, pos int) (VideoNode, error) {
	var videoNode VideoNode
	if pos < 0 || pos > len(tl.VideoNodes) {
		return videoNode, fmt.Errorf("Insertion position is invalid")
	}

	videoNode = createVideoNode(rid, start, end)
	tl.VideoNodes = append(tl.VideoNodes, videoNode)
	return videoNode, nil

}

func (tl *Timeline) Delete(pos int) error {
	if pos < 0 || pos > len(tl.VideoNodes) {
		return fmt.Errorf("Insertion position is invalid")
	}
	tl.VideoNodes = append(tl.VideoNodes[:pos], tl.VideoNodes[pos+1:]...)
	return nil
}

func (tl *Timeline) Split(eventType string, pos int, start, end float64) ([]VideoNode, error) {
	nodes := []VideoNode{}
	if pos < 0 || pos > len(tl.VideoNodes) {
		return nodes, fmt.Errorf("Split position is invalid")
	}

	splitNode := tl.VideoNodes[pos]

	switch eventType {
	case EVT_SLICE_CUT:
		if end > splitNode.Start && end+0.1 < splitNode.End {
			nodes = append(nodes, createVideoNode(splitNode.RID, start, end), createVideoNode(splitNode.RID, end+0.1, splitNode.End))
		}
	case EVT_INTERVAL_CUT:
		if start-0.1 > splitNode.Start && end+0.1 < splitNode.End {
			nodes = append(nodes, createVideoNode(splitNode.RID, splitNode.Start, start-0.1), createVideoNode(splitNode.RID, start, end),
				createVideoNode(splitNode.RID, end+0.1, splitNode.End))
		}
	}

	if len(nodes) <= 0 {
		return nodes, fmt.Errorf("invalid cut range")
	}
	tl.VideoNodes = append(tl.VideoNodes[:pos], append(nodes, tl.VideoNodes[pos+1:]...)...)
	return nodes, nil
}

func (tl *Timeline) MergeClipsQuery(opts *ProcessingOpts) (string, error) {
	if tl.VideoNodes == nil || len(tl.VideoNodes) == 0 {
		return "", fmt.Errorf("no timeline exists")
	}

	var query strings.Builder
	var pos int
	ridToPos := map[string]int{}

	query.WriteString(" -filter_complex \"")
	for i, videoNode := range tl.VideoNodes {
		if pos, ok := ridToPos[videoNode.RID]; ok {
			query.WriteString(fmt.Sprintf("[%d:v]trim=start=%f:end=%f,setpts=PTS-STARTPTS,scale=%s[v%d];", pos, videoNode.Start, videoNode.End, opts.Resolution, i))
			continue
		}
		ridToPos[videoNode.RID] = pos
		query.WriteString(fmt.Sprintf("[%d:v]trim=start=%f:end=%f,setpts=PTS-STARTPTS,scale=%s[v%d];", pos, videoNode.Start, videoNode.End, opts.Resolution, i))
		pos += 1
	}

	for i := range tl.VideoNodes {
		query.WriteString(fmt.Sprintf("[v%d]", i))
	}

	query.WriteString(fmt.Sprintf("concat=n=%d:v=1:a=0[out]\" -map \"[out]\"", len(tl.VideoNodes)))
	return query.String(), nil
}

func (tl *Timeline) InputArgs() (string, error) {
	var pos int
	var query strings.Builder
	inputToPos := map[string]int{}
	for _, videoNode := range tl.VideoNodes {
		if _, ok := inputToPos[videoNode.RID]; ok {
			continue
		}
		inputToPos[videoNode.RID] = pos
		pos += 1
		query.WriteString(fmt.Sprintf(" -i \"%s\"", videoNode.RID))
	}
	return query.String(), nil
}

func (tl *Timeline) OutputArgs(opts *ProcessingOpts) (string, error) {
	if tl.VideoNodes == nil || len(tl.VideoNodes) == 0 {
		return "", fmt.Errorf("no timeline exists")
	}

	return fmt.Sprintf(" -c:v %s -crf %s -preset %s %s", opts.Codec, opts.CRF, opts.Preset, GetFullOutputPath(opts)), nil
}

// CreateProxyFile: creates a copy file from the original to preserve original and work with the given video clip
func CreateProxyFileCMD(inputFilePath, outputFilePath string) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", inputFilePath, // input
		"-codec", "copy",
		"-strict", "experimental",
		outputFilePath)

}

// GenerateEditThumbnail: generate a thumbnail from a video
func GenerateEditThumb(inputFilePath string, outputFilePath string, opts ThumbnailOpts) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", inputFilePath, // input file
		"-vframes", "1", // pick 1 frame from the video
		"-s", "256x256", // scale of the video frame
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
		return "", "", fmt.Errorf("Invalid file format, %s does not contain a file extension", fileName)
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

// GetFullOutputPath: gets the full export path of the video project
func GetFullOutputPath(opts *ProcessingOpts) string {
	return path.Join(opts.OutputPath, opts.Filename+opts.VideoFormat)
}

// ValidateProcessingOpts: validates the processing options of a video project
func ValidateProcessingOpts(opts *ProcessingOpts) error {
	if opts.Filename == "" {
		return fmt.Errorf("filename must be provided")
	}
	if opts.OutputPath == "" {
		return fmt.Errorf("export path was not selected")
	}

	if opts.Codec == "" {
		return fmt.Errorf("codec was not specified")
	}

	if !opts.isCodecCompatible() {
		return fmt.Errorf("codec %s is not compatible", opts.Codec)
	}

	for _, extension := range getValidVideoExtensions() {
		if strings.Contains(opts.Filename, extension) {
			return fmt.Errorf(fmt.Sprintf("invalid filename %s", opts.Filename))
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
