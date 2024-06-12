// query.go: implements high order ffmpeg queries (lossless cut, clip concatenation)
package ffmpegbuilder

import (
	"github.com/k1nho/gahara/internal/video"
)

func CheckVideoDuration(userOpts video.ProcessingOpts) (string, error) {
	input := GetFullInputPath(userOpts)

	query, err := NewDefaultFFmpegBuilder().WithInputs(input).
		WithNullOutput().WithVerbose("").BuildQuery()
	if err != nil {
		return "", err
	}
	return query, nil
}

// CreateProxyFileQuery: creates a proxy file for a video
func CreateProxyFileQuery(userOpts video.ProcessingOpts, format string) (string, error) {
	input := GetFullInputPath(userOpts)
	userOpts.VideoFormat = format
	output := GetFullOutputPath(userOpts)

	querybuilder := NewDefaultFFmpegBuilder().WithInputs(input).WithCodec("copy").
		WithOutputs(output)

	if err := querybuilder.validateProxyFileCreationQuery(); err != nil {
		return "", err
	}

	query, err := querybuilder.BuildQuery()
	if err != nil {
		return "", err
	}

	return query, nil
}

// CreateThumbnailQuery: generates a thumbnail taking the 1 frame of a video
func CreateThumbnailQuery(userOpts video.ProcessingOpts, format string) (string, error) {
	input := GetFullInputPath(userOpts)
	userOpts.VideoFormat = format
	output := GetFullOutputPath(userOpts)

	query, err := NewDefaultFFmpegBuilder().WithInputs(input).WithScale(userOpts.Resolution).
		WithVideoFrames("1").WithOutputs(output).BuildQuery()
	if err != nil {
		return "", err
	}

	return query, nil
}

// MergeClipsQuery: returns the query to concatenate a series of video nodes
func MergeClipsQuery(videoNodes []video.VideoNode, userOpts video.ProcessingOpts) (string, error) {
	querybuilder := NewDefaultFFmpegBuilder().WithInputs(ExtractInputs(videoNodes)...).
		WithPreset(userOpts.Preset).WithCRF(userOpts.CRF).WithVideoCodec(userOpts.Codec).
		WithFScale(userOpts.Resolution).WithOutputs(GetFullOutputPath(userOpts))

	concatFilterQuery, err := querybuilder.ConcatFilter(videoNodes)
	if err != nil {
		return "", err
	}

	querybuilder.ComplexFilterGraph = append(querybuilder.ComplexFilterGraph, concatFilterQuery)
	if err := querybuilder.validateMergeQuery(); err != nil {
		return "", err
	}

	query, err := querybuilder.BuildQuery()
	if err != nil {
		return "", err
	}
	return query, nil
}

// LosslessCutQuery: returns the query string to make a lossless cut of a video node
func LosslessCutQuery(videoNode video.VideoNode, userOpts video.ProcessingOpts) (string, error) {
	// overwrite filename, if it was passed by default lossy opts
	userOpts.Filename = videoNode.Name

	querybuilder := NewDefaultFFmpegBuilder().WithInputs(videoNode.RID).WithInputStartTime(videoNode.Start).
		WithOutputDuration(videoNode.End - videoNode.Start).WithCodec("copy").WithAvoidNegativeTS("make_zero").
		WithMovFlags("+faststart").WithOutputs(GetFullOutputPath(userOpts))

	if err := querybuilder.validateLosslessCutQuery(); err != nil {
		return "", err
	}

	query, err := querybuilder.BuildQuery()
	if err != nil {
		return "", err
	}
	return query, nil
}
