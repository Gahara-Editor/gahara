package timeline

import (
	"fmt"
	"slices"

	"github.com/k1nho/gahara/internal/audio"
	"github.com/k1nho/gahara/internal/placeholder"
	"github.com/k1nho/gahara/internal/text"
	"github.com/k1nho/gahara/internal/video"
)

func NewTimeline() Timeline {
	t := Timeline{Nodes: make([][]TimelineNode, 0)}
	return t
}

func CreateNode(nodeType string, rid string, name string, start, end float64) TimelineNode {
	switch nodeType {
	case video.NODE_VIDEO:
		if name == "" {
			name = "VIDEO"
		}
		return video.CreateNode(rid, name, start, end)
	case audio.NODE_AUDIO:
		if name == "" {
			name = "AUDIO"
		}
		return audio.CreateNode(rid, name, start, end)
	case text.NODE_TEXT:
		if name == "" {
			name = "TEXT"
		}

	}

	return placeholder.CreateNode(rid, name, start, end)
}

func (t *Timeline) AddTrack() {
	t.Nodes = append(t.Nodes, []TimelineNode{})
}

func (t *Timeline) RemoveTrack(tid int) error {
	if tid < 0 || tid >= len(t.Nodes) {
		return fmt.Errorf("invalid track index")
	}
	t.Nodes = append(t.Nodes[:tid], t.Nodes[tid+1:]...)
	return nil
}

func (t *Timeline) Insert(tid int, pos int, node TimelineNode) (TimelineNode, error) {
	if len(t.Nodes) == 0 {
		t.AddTrack()
	}
	if !(tid >= 0 && tid < len(t.Nodes) && pos >= 0 && pos <= len(t.Nodes[tid])) {
		return node, fmt.Errorf("insertion position %d is invalid, track does not exists", tid)
	}
	t.Nodes[tid] = slices.Insert(t.Nodes[tid], pos, node)
	return node, nil
}

func (t *Timeline) Delete(tid int, pos int) error {
	if !t.inBounds(tid, pos) {
		return fmt.Errorf("delete position is invalid")
	}
	if len(t.Nodes[tid]) == 0 {
		return fmt.Errorf("there are no video clips to delete in track")
	}
	t.Nodes[tid] = slices.Delete(t.Nodes[tid], pos, pos+1)
	return nil
}

func (t *Timeline) Split(tid int, pos int, eventType string, start, end float64) ([]TimelineNode, error) {
	nodes := []TimelineNode{}
	if !t.inBounds(tid, pos) {
		return nodes, fmt.Errorf("delete position is invalid")
	}
	if len(t.Nodes[tid]) == 0 {
		return nodes, fmt.Errorf("there are no video clips to delete in track")
	}

	splitNode := t.Nodes[tid][pos]

	switch eventType {
	case video.EVT_SLICE_CUT:
		if end > splitNode.Start() && end+0.1 < splitNode.End() {
			nodes = append(nodes, CreateNode(splitNode.Type(), splitNode.RID(), splitNode.Name(), start, end),
				CreateNode(splitNode.Type(), splitNode.RID(), splitNode.Name(), end+0.1, splitNode.End()))
		}
	case video.EVT_INTERVAL_CUT:
		if start-0.1 > splitNode.Start() && end+0.1 < splitNode.End() {
			nodes = append(nodes, CreateNode(splitNode.Type(), splitNode.RID(), splitNode.Name(), splitNode.Start(), start-0.1),
				CreateNode(splitNode.Type(), splitNode.RID(), splitNode.Name(), start, end),
				CreateNode(splitNode.Type(), splitNode.RID(), splitNode.Name(), end+0.1, splitNode.End()))
		}
	}

	if len(nodes) <= 0 {
		return nodes, fmt.Errorf("invalid cut range")
	}
	t.Nodes[tid] = append(t.Nodes[tid][:pos], append(nodes, t.Nodes[tid][pos+1:]...)...)
	return nodes, nil
}

func (t *Timeline) DeleteRIDReferences(rid string) error {
	if t.Nodes == nil {
		return fmt.Errorf("no timeline exists")
	}

	for i := range t.Nodes {
		t.Nodes[i] = slices.DeleteFunc(t.Nodes[i], func(node TimelineNode) bool {
			return node.RID() == rid
		})

	}
	return nil
}

func (t *Timeline) inBounds(tid int, pos int) bool {
	return tid >= 0 && tid < len(t.Nodes) && pos >= 0 && pos < len(t.Nodes[tid])
}

func (t *Timeline) RenameVideoNode(tid int, pos int, name string) error {
	if !t.inBounds(tid, pos) {
		return fmt.Errorf("clip position invalid")
	}
	if len(t.Nodes[tid]) == 0 {
		return fmt.Errorf("there are no video clips to rename in track")
	}
	t.Nodes[tid][pos].Rename(name)
	return nil
}

func (t *Timeline) ToggleLossless(tid int, pos int) error {
	if !t.inBounds(tid, pos) {
		return fmt.Errorf("clip position invalid (%d)", pos)
	}
	if len(t.Nodes[tid]) == 0 {
		return fmt.Errorf("there are no video clips to rename in track")
	}

	if t.Nodes[tid][pos].Type() == video.NODE_VIDEO {
		t.Nodes[tid][pos].(*video.VideoNode).LosslessExport = !t.Nodes[tid][pos].(*video.VideoNode).LosslessExport
	}
	return nil
}

func (t *Timeline) MarkAllLossless(tid int) error {
	if !t.inBounds(tid, 0) {
		return fmt.Errorf("clip position invalid")
	}
	if len(t.Nodes[tid]) == 0 {
		return fmt.Errorf("there are no video clips to rename in track")
	}

	for i := range t.Nodes[tid] {
		if t.Nodes[tid][i].Type() == video.NODE_VIDEO {
			t.Nodes[tid][i].(*video.VideoNode).LosslessExport = true
		}
	}

	return nil
}

func (t *Timeline) UnmarkAllLossless(tid int) error {
	if !t.inBounds(tid, 0) {
		return fmt.Errorf("clip position invalid")
	}
	if len(t.Nodes[tid]) == 0 {
		return fmt.Errorf("there are no video clips to rename in track")
	}

	for i := range t.Nodes[tid] {
		if t.Nodes[tid][i].Type() == video.NODE_VIDEO {
			t.Nodes[tid][i].(*video.VideoNode).LosslessExport = false
		}
	}

	return nil
}
