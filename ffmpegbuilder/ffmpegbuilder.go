package ffmpegbuilder

import (
	"fmt"
	"strings"

	"github.com/k1nho/gahara/internal/video"
)

type FFmpegBuilder struct {
	PreInputParams     PreInputParams
	Inputs             []string
	FilterGraphParams  FilterGraphParams
	ComplexFilterGraph []string
	OutputParams       OutputParams
	Outputs            []string
}

// PreInputParams: all the parameters for input
type PreInputParams struct {
	// VerboseMode: -v in fffmpeg (sets the log level quiet, panic, fatal, error, warning, info, verbose, debug, trace)
	VerboseMode string
	// StatsPeriod: -stats_period in ffmpeg (sets the period at which encoding progress/statistics are updated)
	StatsPeriod string
	// Progress: -progress in ffmpeg (sends the progress information to a url, it can be also stdout->pipe:1, stderr->pipe:2)
	Progress string
	// HideBanner: -hide_banner in ffmpeg (hides the default CLI banner)
	HideBanner bool
	// StartTime: -ss in ffmpeg (seek time, or start time of a video)
	StartTime float64
}

type FilterGraphParams struct {
	// Scale: -s in ffmpeg, the scale of the output (1920x1080)
	Scale string
}

// OutputParams: all the parameters for output
type OutputParams struct {
	// Codec: -c in ffmpeg
	Codec string
	// VideoCodec: -c:v in ffmpeg (H264, H265, copy)
	VideoCodec string
	// AudioCodec: -c:a in ffmpeg
	AudioCodec string
	// Duration: -t in ffmpeg, represents how long should the video last from a StartTime (00:00:20, 42.37)
	Duration float64
	// StopTime: -to in ffmpeg, represents when should the the video stop reading or writing(00:00:20, 42.37),
	// if Duration is specified it will Duration will have priority
	StopTime float64
	// CopyTS: -copyts in ffmpeg, keeps the original timestamps
	CopyTS bool
	// VideoFrames: -frames:v in ffmpeg, the number of video frames to output
	VideoFrames string
	// Scale: -s in ffmpeg, the scale of the output (1920x1080)
	Scale string
	//  CRF: -crf in ffmpeg, the constant rate factor (the lower the better)
	CRF string
	// Preset: -preset in ffmpeg, encoding speed to compression ratio (slow, medium, fast)
	Preset string
	// AvoidNegativeTS: -avoid_negative_ts in ffmpeg, avoids negative timestamps (make_zero: first ts 0, make_non_negative, disabled)
	AvoidNegativeTS string
	// MovFlags: -movflags in ffmpeg, mov, mp4, and ismv support fragmentation. The metadata about all packets is stored in one location,
	// but it can be moved at the start for better playback (+faststart)
	MovFlags string
}

func NewDefaultFFmpegBuilder() *FFmpegBuilder {
	return &FFmpegBuilder{
		PreInputParams:     NewDefaultPreInputParams(),
		Inputs:             []string{},
		ComplexFilterGraph: []string{},
		FilterGraphParams:  NewDefaultFiltergraphParams(),
		OutputParams:       NewDefaultOutputParams(),
		Outputs:            []string{},
	}
}

func NewDefaultOutputParams() OutputParams {
	return OutputParams{}
}

func NewDefaultPreInputParams() PreInputParams {
	return PreInputParams{
		VerboseMode: "quiet",
		StatsPeriod: "5s",
		Progress:    "pipe:2",
		HideBanner:  true,
	}
}

func NewDefaultFiltergraphParams() FilterGraphParams {
	return FilterGraphParams{}
}

func (f *FFmpegBuilder) WithInputs(inputs ...string) *FFmpegBuilder {
	f.Inputs = append(f.Inputs, inputs...)
	return f
}

func (f *FFmpegBuilder) WithOutputs(outputs ...string) *FFmpegBuilder {
	f.Outputs = append(f.Outputs, outputs...)
	return f
}

func (f *FFmpegBuilder) WithInputStartTime(startTime float64) *FFmpegBuilder {
	f.PreInputParams.StartTime = startTime
	return f
}

func (f *FFmpegBuilder) WithStatsPeriod(statsPeriod string) *FFmpegBuilder {
	f.PreInputParams.StatsPeriod = statsPeriod
	return f
}

func (f *FFmpegBuilder) WithCodec(codec string) *FFmpegBuilder {
	f.OutputParams.Codec = codec
	return f
}

func (f *FFmpegBuilder) WithVideoCodec(videoCodec string) *FFmpegBuilder {
	f.OutputParams.VideoCodec = videoCodec
	return f
}

func (f *FFmpegBuilder) WithAudioCodec(audioCodec string) *FFmpegBuilder {
	f.OutputParams.AudioCodec = audioCodec
	return f
}

func (f *FFmpegBuilder) WithOutputStopTime(stopTime float64) *FFmpegBuilder {
	f.OutputParams.StopTime = stopTime
	return f
}

func (f *FFmpegBuilder) WithOutputDuration(duration float64) *FFmpegBuilder {
	f.OutputParams.Duration = duration
	return f
}

func (f *FFmpegBuilder) WithMovFlags(flag string) *FFmpegBuilder {
	f.OutputParams.MovFlags = flag
	return f
}

func (f *FFmpegBuilder) WithCRF(crf string) *FFmpegBuilder {
	f.OutputParams.CRF = crf
	return f
}

func (f *FFmpegBuilder) WithPreset(preset string) *FFmpegBuilder {
	f.OutputParams.Preset = preset
	return f
}

func (f *FFmpegBuilder) WithCopyTS() *FFmpegBuilder {
	f.OutputParams.CopyTS = true
	return f
}

func (f *FFmpegBuilder) WithAvoidNegativeTS(mode string) *FFmpegBuilder {
	f.OutputParams.AvoidNegativeTS = mode
	return f
}

func (f *FFmpegBuilder) WithVideoFrames(nFrames string) *FFmpegBuilder {
	f.OutputParams.VideoFrames = nFrames
	return f
}

func (f *FFmpegBuilder) WithScale(scale string) *FFmpegBuilder {
	f.OutputParams.Scale = scale
	return f
}

// WithFScale: sets the resolution to be used wiithin the filtergraph
func (f *FFmpegBuilder) WithFScale(scale string) *FFmpegBuilder {
	f.FilterGraphParams.Scale = scale
	return f
}

// ConcatFilter: returns a concatenation query to be used in filter complex, given video nodes
func (f *FFmpegBuilder) ConcatFilter(videoNodes []video.VideoNode) (string, error) {
	if len(videoNodes) == 0 {
		return "", fmt.Errorf("no video nodes were provided")
	}

	var concatQuery strings.Builder
	var pos int
	ridToPos := map[string]int{}

	concatQuery.WriteString("\"")
	for i, videoNode := range videoNodes {
		if pos, ok := ridToPos[videoNode.RID]; ok {
			concatQuery.WriteString(fmt.Sprintf("[%d:v]trim=start=%.4f:end=%.4f,setpts=PTS-STARTPTS,scale=%s[v%d];", pos, videoNode.Start, videoNode.End, f.FilterGraphParams.Scale, i))
			continue
		}
		ridToPos[videoNode.RID] = pos
		concatQuery.WriteString(fmt.Sprintf("[%d:v]trim=start=%.4f:end=%.4f,setpts=PTS-STARTPTS,scale=%s[v%d];", pos, videoNode.Start, videoNode.End, f.FilterGraphParams.Scale, i))
		pos += 1
	}

	for i := range videoNodes {
		concatQuery.WriteString(fmt.Sprintf("[v%d]", i))
	}

	concatQuery.WriteString(fmt.Sprintf("concat=n=%d:v=1:a=0[out]\" -map \"[out]\"", len(videoNodes)))
	return concatQuery.String(), nil
}

// BuildQuery: returns the ffmpeg query with all the parameters given
func (f *FFmpegBuilder) BuildQuery() (string, error) {
	var cmd strings.Builder
	cmd.WriteString("ffmpeg ")

	if f.PreInputParams.HideBanner {
		cmd.WriteString("-hide_banner ")
	}
	if f.PreInputParams.VerboseMode != "" {
		cmd.WriteString("-v ")
		cmd.WriteString(f.PreInputParams.VerboseMode)
		cmd.WriteString(" ")
	}
	if f.PreInputParams.StatsPeriod != "" {
		cmd.WriteString("-stats_period ")
		cmd.WriteString(f.PreInputParams.StatsPeriod)
		cmd.WriteString(" ")
	}
	if f.PreInputParams.Progress != "" {
		cmd.WriteString("-progress ")
		cmd.WriteString(f.PreInputParams.Progress)
		cmd.WriteString(" ")
	}

	if f.PreInputParams.StartTime != 0 {
		cmd.WriteString(fmt.Sprintf("-ss %.4f ", f.PreInputParams.StartTime))
	}

	// Append inputs
	dups := make(map[string]struct{})
	for _, input := range f.Inputs {
		if _, ok := dups[input]; ok {
			continue
		}
		cmd.WriteString("-i ")
		cmd.WriteString(fmt.Sprintf("\"%s\"", input))
		cmd.WriteString(" ")
		dups[input] = struct{}{}
	}

	// Append complex filter graph
	if len(f.ComplexFilterGraph) > 0 {
		cmd.WriteString("-filter_complex ")
		for _, filterGraph := range f.ComplexFilterGraph {
			cmd.WriteString(filterGraph)
			cmd.WriteString(" ")
		}
	}

	// Append output parameters
	if f.OutputParams.Duration != 0 {
		cmd.WriteString(fmt.Sprintf("-t %.4f ", f.OutputParams.Duration))
	}
	if f.OutputParams.StopTime != 0 {
		cmd.WriteString(fmt.Sprintf("-to %.4f ", f.OutputParams.StopTime))
	}

	if f.OutputParams.AvoidNegativeTS != "" {
		cmd.WriteString(fmt.Sprintf("-avoid_negative_ts %s ", f.OutputParams.AvoidNegativeTS))
	}

	if f.OutputParams.Codec != "" {
		cmd.WriteString("-c ")
		cmd.WriteString(f.OutputParams.Codec)
		cmd.WriteString(" ")
	}
	if f.OutputParams.VideoCodec != "" {
		cmd.WriteString("-c:v ")
		cmd.WriteString(f.OutputParams.VideoCodec)
		cmd.WriteString(" ")
	}
	if f.OutputParams.AudioCodec != "" {
		cmd.WriteString("-c:a ")
		cmd.WriteString(f.OutputParams.AudioCodec)
		cmd.WriteString(" ")
	}

	if f.OutputParams.MovFlags != "" {
		cmd.WriteString(fmt.Sprintf("-movflags '%s' ", f.OutputParams.MovFlags))
	}

	if f.OutputParams.CRF != "" {
		cmd.WriteString("-crf ")
		cmd.WriteString(f.OutputParams.CRF)
		cmd.WriteString(" ")
	}
	if f.OutputParams.Preset != "" {
		cmd.WriteString("-preset ")
		cmd.WriteString(f.OutputParams.Preset)
		cmd.WriteString(" ")
	}

	if f.OutputParams.CopyTS {
		cmd.WriteString("-copyts ")
	}

	if f.OutputParams.VideoFrames != "" {
		cmd.WriteString("-frames:v ")
		cmd.WriteString(f.OutputParams.VideoFrames)
		cmd.WriteString(" ")
	}
	if f.OutputParams.Scale != "" {
		cmd.WriteString("-s ")
		cmd.WriteString(f.OutputParams.Scale)
		cmd.WriteString(" ")
	}

	// Append outputs
	for _, output := range f.Outputs {
		cmd.WriteString(fmt.Sprintf("\"%s\"", output))
		cmd.WriteString(" ")
	}

	return cmd.String(), nil
}
