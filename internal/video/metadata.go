// metadata.go contains all the data structures, and operations that preserve video editing metadata
// for example: if the user cuts the video, metadata will contain this in order to be apply on rendering step
package video

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	// EVT_INTERVAL_CUT: interval cut event
	EVT_INTERVAL_CUT = "intervalCut"
	// EVT_SLICE_CUT: slice cut event
	EVT_SLICE_CUT = "sliceCut"
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
