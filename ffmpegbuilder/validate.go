package ffmpegbuilder

import "fmt"

func (f *FFmpegBuilder) validateMergeQuery() error {
	if len(f.Inputs) == 0 {
		return fmt.Errorf("no input stream was provided")
	}
	if len(f.Outputs) == 0 {
		return fmt.Errorf("no output stream was provided")
	}
	if f.FilterGraphParams.Scale == "" {
		return fmt.Errorf("no resolution was provided. For this operations clips need to have the same resolution")
	}
	if f.OutputParams.VideoCodec == "" {
		return fmt.Errorf("no codec was provided")
	}

	if f.OutputParams.Preset == "" {
		return fmt.Errorf("no preset was provided")
	}
	if f.OutputParams.CRF == "" {
		return fmt.Errorf("no constant rate factor was provided")
	}
	return nil
}

func (f *FFmpegBuilder) validateLosslessCutQuery() error {
	if len(f.Inputs) != 1 {
		return fmt.Errorf("no input stream(s) provided")
	}
	if len(f.Outputs) != 1 {
		return fmt.Errorf("no output stream(s) provided")
	}

	if f.OutputParams.Duration < 0 {
		return fmt.Errorf("specified duration of the clip is negative")
	}
	return nil
}

func (f *FFmpegBuilder) validateProxyFileCreationQuery() error {
	if len(f.Inputs) != 1 {
		return fmt.Errorf("no input stream(s) provided")
	}
	if len(f.Outputs) != 1 {
		return fmt.Errorf("no output stream(s) provided")
	}

	if f.OutputParams.Codec != "copy" {
		return fmt.Errorf("codec must be copy. No re-encoding needed for proxy")
	}
	return nil
}
