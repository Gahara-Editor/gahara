package ffmpegbuilder

import (
	"path"

	"github.com/k1nho/gahara/internal/audio"
	"github.com/k1nho/gahara/internal/timeline"
	"github.com/k1nho/gahara/internal/video"
)

// extractInputs: extract input params from video nodes
func ExtractInputs(videoNodes []video.VideoNode) []string {
	var inputs []string
	for _, videoNode := range videoNodes {
		inputs = append(inputs, videoNode.VideoRID)
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

func ExtractVideoNodes(track []timeline.TimelineNode) []video.VideoNode {
	videoNodes := make([]video.VideoNode, 0)
	for _, node := range track {
		if node.Type() == video.NODE_VIDEO {
			if videoNode, ok := node.(*video.VideoNode); ok {
				videoNodes = append(videoNodes, *videoNode)
			}
		}
	}
	return videoNodes
}

func ExtractAudioNodes(track []timeline.TimelineNode) []audio.AudioNode {
	audioNodes := make([]audio.AudioNode, 0)
	for _, node := range track {
		if node.Type() == audio.NODE_AUDIO {
			if audioNode, ok := node.(*audio.AudioNode); ok {
				audioNodes = append(audioNodes, *audioNode)
			}
		}
	}
	return audioNodes
}
