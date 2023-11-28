// package video contains all the video editing commands that interface with ffmpeg
package video

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/k1nho/gahara/internal/utils"
)

type ThumbnailOpts struct {
	StartTime string
	Scale     string
}

// CreateProxyFile: creates a copy file from the original to preserve original and work with the given video clip
func CreateProxyFileCMD(path, pathProxyFile string) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-i", path, // input
		"-codec", "copy",
		"-strict", "experimental",
		pathProxyFile)

}

// CutVideoInterval: given a file it returns the video with the interval (start,end) removed
func CutVideoInterval(file string, start, end time.Duration) *exec.Cmd {
	timeStart, timeEnd := utils.FormatTime(start), utils.FormatTime(end)

	return exec.Command("ffmpeg",
		"-i", file, // input file
		"-ss", timeStart, // start of the cut
		"-to", timeEnd, // end of the cut
		"-codec copy", // avoid re-encoding only need to trim the video (lossless cut)
		file,          // output file will be the same
	)
}

func AddVideoClipInInterval(mainVideo, clipVideo string, start, end time.Duration) *exec.Cmd {
	timeStart, timeEnd := utils.FormatTime(start), utils.FormatTime(end)
	return exec.Command("ffmpeg",
		"-i", mainVideo, // input main video
		"-i", clipVideo, // input clip video
		"-filter_complex", // use complex filtergraph (multiple input files or output is different than the input)
		fmt.Sprintf("[0:v]trim=duration=%s[s1];[0:v]trim=start=%s:end=%s[s2];[s1][1:v][s2]concat=n=3:v=1:a=0", timeStart, timeStart, timeEnd),
		"-c:v", "libx264",
		"c:a", "acc",
		"-strict", "experimental",
		"aci.mov",
	)
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
