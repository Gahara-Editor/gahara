// video.go contains all the data structures, and operations for video editing
package video

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

var ValidVideoExtensions = [16]string{".mov", ".mkv", ".mp4", ".m4v", ".mp4v", ".avi", ".wmv", ".flv", ".h264", ".hevc", ".3gp", ".3g2", ".ogv", ".264", ".265", ".webm"}

const (
	Epsilon = 1e-6
	// EVT_INTERVAL_CUT: interval cut event
	EVT_INTERVAL_CUT = "intervalCut"
	// EVT_SLICE_CUT: slice cut event
	EVT_SLICE_CUT = "sliceCut"
	// SCALE_1920X1080: Resolution 1920x1080
	SCALE_1920X1080 = "1920:1080"
	// SCALE_1920X1080: Resolution 2560x1440
	SCALE_2560x1440 = "2560x1440"
	// SCALE_3840X2160: Resolution 3840x2160
	SCALE_3840X2160 = "3840:2160"
	// CODEC_H264: codec H.264
	CODEC_H264 = "libx264"
	// CODEC_H264_RGB: codec H.264 RGB
	CODEC_H264_RGB = "libx264rgb"
	// CODEC_COPY: codec copy for bitstream copy (skips re-encoding the video)
	CODEC_COPY = "copy"
	// CRF_23: constant rate factor 23, lower number usually means better quality
	CRF_23 = "23"
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
	StartTime string
	Scale     string
}

type ProcessingOpts struct {
	// Resolution (1920x1080, 2560x1440, 3840x2160)
	Resolution string
	// Codec: video codec (libx264, libx265)
	Codec string
	// CRF: Constant Rate Factor
	CRF string
	// Preset: Encoding speed to compression ratio
	Preset string
	// OutputPath: The file path of the output video
	OutputPath string
}

func NewDefaultProcessingOpts() *ProcessingOpts {
	return &ProcessingOpts{
		OutputPath: "output.mp4",
		Resolution: SCALE_1920X1080,
		Codec:      CODEC_H264_RGB,
		CRF:        CRF_18,
		Preset:     PRESET_MEDIUM,
	}
}

func (o *ProcessingOpts) WithResolution(res string) {
	o.Resolution = res
}

func (o *ProcessingOpts) WithCodec(codec string) {
	o.Codec = codec
}

func (o *ProcessingOpts) WithCRF(crf string) {
	o.CRF = crf
}

func (o *ProcessingOpts) WithPreset(preset string) {
	o.Preset = preset
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

	return fmt.Sprintf(" -c:v %s -crf %s -preset %s %s", opts.Codec, opts.CRF, opts.Preset, opts.OutputPath), nil
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
