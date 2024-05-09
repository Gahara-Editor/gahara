package ffmpegbuilder

import (
	"path"

	"github.com/k1nho/gahara/internal/video"
)

// extractInputs: extract input params from video nodes
func ExtractInputs(videoNodes []video.VideoNode) []string {
	var inputs []string
	for _, videoNode := range videoNodes {
		inputs = append(inputs, videoNode.RID)
	}
	return inputs
}

// getFullOutputPath: gets the full export path of the video
func GetFullOutputPath(opts video.ProcessingOpts) string {
	return path.Join(opts.OutputPath, opts.Filename+opts.VideoFormat)
}

// getFullOutputPath: gets the full input path of the video
func GetFullInputPath(opts video.ProcessingOpts) string {
	return path.Join(opts.InputPath, opts.Filename+opts.VideoFormat)
}
