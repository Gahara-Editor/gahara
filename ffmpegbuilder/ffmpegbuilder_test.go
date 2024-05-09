package ffmpegbuilder

import (
	"fmt"
	"testing"

	"github.com/k1nho/gahara/internal/video"
)

func mockTl() *video.Timeline {
	return &video.Timeline{VideoNodes: []video.VideoNode{
		{RID: "root1", ID: "1", Name: "input1", Start: 20.1, End: 25.2, LosslessExport: true},
		{RID: "root1", ID: "2", Name: "input2", Start: 1.12, End: 10.2, LosslessExport: true},
		{RID: "root2", ID: "3", Name: "input3", Start: 12.2, End: 21.2, LosslessExport: true},
		{RID: "root3", ID: "4", Name: "input4", Start: 69.112, End: 80.23, LosslessExport: true},
	}}

}

func TestFFmpegBuilder(t *testing.T) {
	t.Parallel()

	t.Run("format conversion query", func(t *testing.T) {
		expectedQuery := "ffmpeg -hide_banner -v quiet -stats_period 5s -progress pipe:2 -i \"myinput.mp4\" \"myoutput.mov\" "
		query, err := NewDefaultFFmpegBuilder().WithInputs("myinput.mp4").WithOutputs("myoutput.mov").BuildQuery()
		if err != nil {
			t.Fatal(err)
		}
		if query != expectedQuery {
			t.Errorf("\ngot: %s\nexp: %s", query, expectedQuery)
		}
	})

	t.Run("generate proxy file query", func(t *testing.T) {
		expectedQuery := "ffmpeg -hide_banner -v quiet -stats_period 5s -progress pipe:2 -i \"inputpath/input.mp4\" -c copy \"outputpath/input.mov\" "
		query, err := CreateProxyFileQuery(video.ProcessingOpts{
			Filename:    "input",
			InputPath:   "inputpath",
			OutputPath:  "outputpath",
			VideoFormat: ".mp4",
		}, ".mov")
		if err != nil {
			t.Fatal(err)
		}
		if query != expectedQuery {
			t.Errorf("\ngot: %s\nexp: %s", query, expectedQuery)
		}
	})

	t.Run("generate thumbnail query", func(t *testing.T) {
		expectedQuery := "ffmpeg -hide_banner -v quiet -stats_period 5s -progress pipe:2 -i \"inputpath/input.mp4\" -frames:v 1 \"outputpath/input.png\" "
		query, err := CreateThumbnailQuery(video.ProcessingOpts{
			Filename:    "input",
			InputPath:   "inputpath",
			OutputPath:  "outputpath",
			VideoFormat: ".mp4",
		}, ".png")
		if err != nil {
			t.Fatal(err)
		}
		if query != expectedQuery {
			t.Errorf("\ngot: %s\nexp: %s", query, expectedQuery)
		}
	})

	t.Run("concat filter query", func(t *testing.T) {
		expectedQuery := "ffmpeg -hide_banner -v quiet -stats_period 5s -progress pipe:2 -i \"root1\" -i \"root2\" -i \"root3\" -filter_complex \"[0:v]trim=start=20.1000:end=25.2000,setpts=PTS-STARTPTS,scale=1920x1080[v0];[0:v]trim=start=1.1200:end=10.2000,setpts=PTS-STARTPTS,scale=1920x1080[v1];[1:v]trim=start=12.2000:end=21.2000,setpts=PTS-STARTPTS,scale=1920x1080[v2];[2:v]trim=start=69.1120:end=80.2300,setpts=PTS-STARTPTS,scale=1920x1080[v3];[v0][v1][v2][v3]concat=n=4:v=1:a=0[out]\" -map \"[out]\" -c:v libx264 -crf 18 -preset medium \"outputpath/myvideo.mp4\" "

		query, err := MergeClipsQuery(mockTl().VideoNodes, video.ProcessingOpts{
			Resolution:  "1920x1080",
			Codec:       "libx264",
			CRF:         "18",
			Preset:      "medium",
			VideoFormat: ".mp4",
			OutputPath:  "outputpath",
			Filename:    "myvideo",
		})
		if err != nil {
			t.Fatal(err)
		}

		if query != expectedQuery {
			t.Errorf("\ngot: %s\nexp: %s", query, expectedQuery)
		}
	})

	t.Run("lossless cut query", func(t *testing.T) {
		videoNode := video.VideoNode{RID: "root1", Name: "myvideo", Start: 22.2300, End: 28.4321, ID: "1", LosslessExport: true}
		duration := videoNode.End - videoNode.Start

		expectedQuery := fmt.Sprintf("ffmpeg -hide_banner -v quiet -stats_period 5s -progress pipe:2 -ss 22.2300 -i \"root1\" -t %.4f -avoid_negative_ts make_zero -c copy -movflags '+faststart' \"outputpath/myvideo.mp4\" ", duration)

		query, err := LosslessCutQuery(videoNode, video.ProcessingOpts{
			OutputPath:  "outputpath",
			Filename:    videoNode.Name,
			VideoFormat: ".mp4",
		})

		if err != nil {
			t.Fatal(err)
		}
		if query != expectedQuery {
			t.Errorf("\ngot: %s\nexp: %s", query, expectedQuery)
		}
	})
}
