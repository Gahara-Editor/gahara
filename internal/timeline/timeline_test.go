package timeline

import (
	"math"
	"testing"

	"github.com/k1nho/gahara/internal/video"
)

func mockTl() *Timeline {
	return &Timeline{Nodes: [][]TimelineNode{{
		CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
		CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
		CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
		CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
		CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9)}}}
}

func emptyTL() *Timeline {
	t := &Timeline{Nodes: [][]TimelineNode{}}
	t.AddTrack()
	return t
}

func TestTimelineInsert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		element  TimelineNode
		idx      int
		tl       *Timeline
		expected *Timeline
	}{
		{
			name:    "Insert an element in an empty timeline",
			element: CreateNode(video.NODE_VIDEO, "first", "Node", 4.2, 6.9),
			tl:      emptyTL(),
			idx:     0,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "first", "Node", 4.2, 6.9)}},
			},
		},

		{
			name:    "Insert an element in the middle",
			element: CreateNode(video.NODE_VIDEO, "middle", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     2,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "middle", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9)}},
			},
		},
		{
			name:    "Insert an element at the beginning",
			element: CreateNode(video.NODE_VIDEO, "first", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     0,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "first", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name:    "Insert an element at the end",
			element: CreateNode(video.NODE_VIDEO, "last", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     4,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "last", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name:    "Index out out bounds (slice size = 4, idx = -1)",
			element: CreateNode(video.NODE_VIDEO, "fail", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     -1,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name:    "Index out out bounds (slice size = 4, idx = 10)",
			element: CreateNode(video.NODE_VIDEO, "fail", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     6,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.tl.Insert(0, tt.idx, tt.element)
		if tt.idx < 0 || tt.idx > len(tt.tl.Nodes[0]) {
			if err == nil {
				t.Errorf("[%s] index %d is out of bounds and it should have failed", tt.name, tt.idx)
			}
			return
		}

		if err != nil {
			t.Errorf("[%s] failed to insert: %s", tt.name, err.Error())
		}

		if len(tt.tl.Nodes[0]) != len(tt.expected.Nodes[0]) {
			t.Errorf("[%s] the timeline lenghts do not match: (expected: %d, got: %d)", tt.name, len(tt.expected.Nodes[0]), len(tt.tl.Nodes[0]))
		}

		for i := range tt.tl.Nodes[0] {
			a, oka := tt.tl.Nodes[0][i].(*video.VideoNode)
			b, okb := tt.expected.Nodes[0][i].(*video.VideoNode)
			if !oka || !okb {
				t.Errorf("[%s] video node expected, got %s, %s", tt.name, a.Type(), b.Type())
			}
			if a.VideoRID != b.VideoRID &&
				a.VideoName != b.VideoName &&
				math.Abs(a.VideoStart-b.VideoStart) <= 0.01 &&
				math.Abs(a.VideoEnd-b.VideoEnd) <= 0.01 {
				t.Errorf("[%s] the timelines do not have the same order", tt.name)
			}
		}

	}
}

func TestTimelineDelete(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		idx      int
		tl       *Timeline
		expected *Timeline
	}{
		{
			name: "Delete an element in the middle",
			tl:   mockTl(),
			idx:  2,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name: "Delete the first element",
			tl:   mockTl(),
			idx:  0,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name: "Delete the last element",
			tl:   mockTl(),
			idx:  0,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name: "Index out out bounds (slice size = 4, idx = -1)",
			tl:   mockTl(),
			idx:  -1,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
		{
			name: "Index out out bounds (slice size = 4, idx = 10)",
			tl:   mockTl(),
			idx:  10,
			expected: &Timeline{Nodes: [][]TimelineNode{{
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			},
			},
			},
		},
	}

	for _, tt := range tests {
		err := tt.tl.Delete(0, tt.idx)
		if tt.idx < 0 || tt.idx > len(tt.tl.Nodes[0]) {
			if err == nil {
				t.Errorf("[%s] index %d is out of bounds and it should have failed", tt.name, tt.idx)
			}
			return
		}

		if err != nil {
			t.Errorf("[%s] could not perform delete operation: %s", tt.name, err.Error())
		}

		if len(tt.tl.Nodes[0]) != len(tt.expected.Nodes[0]) {
			t.Errorf("[%s] the timeline lenghts do not match: (expected: %d, got: %d)", tt.name, len(tt.expected.Nodes[0]), len(tt.tl.Nodes[0]))
		}

		for i := range tt.tl.Nodes[0] {
			a, oka := tt.tl.Nodes[0][i].(*video.VideoNode)
			b, okb := tt.expected.Nodes[0][i].(*video.VideoNode)
			if !oka || !okb {
				t.Errorf("[%s] video node expected, got %s, %s", tt.name, a.Type(), b.Type())
			}
			if a.VideoRID != b.VideoRID &&
				a.VideoName != b.VideoName &&
				math.Abs(a.VideoStart-b.VideoStart) <= 0.01 &&
				math.Abs(a.VideoEnd-b.VideoEnd) <= 0.01 {
				t.Errorf("[%s] the timelines do not have the same order", tt.name)
			}
		}

	}
}

func TestDeleteRIDByReferences(t *testing.T) {
	t.Run("delete references of a rid", func(t *testing.T) {
		expected := &Timeline{
			Nodes: [][]TimelineNode{{CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
				CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
			},
			},
		}

		tl := &Timeline{Nodes: [][]TimelineNode{{
			CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "2", "Node", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "3", "Node", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "1", "Node", 4.2, 6.9)}}}

		err := tl.DeleteRIDReferences("1")
		if err != nil {
			t.Errorf("failed to perform reference deletion")
		}

		if len(tl.Nodes[0]) != len(expected.Nodes[0]) {
			t.Errorf("the timeline lenghts do not match: (expected: %d, got: %d)", len(expected.Nodes[0]), len(tl.Nodes[0]))
		}

		for i := range tl.Nodes[0] {
			a, oka := tl.Nodes[0][i].(*video.VideoNode)
			b, okb := expected.Nodes[0][i].(*video.VideoNode)
			if !oka || !okb {
				t.Errorf("video node expected, got %s, %s", a.Type(), b.Type())
			}
			if a.VideoRID != b.VideoRID &&
				a.VideoName != b.VideoName &&
				math.Abs(a.VideoStart-b.VideoStart) <= 0.01 &&
				math.Abs(a.VideoEnd-b.VideoEnd) <= 0.01 {
				t.Errorf("the timelines do not have the same order")
			}
		}

	})
}

func TestMarkLossless(t *testing.T) {
	t.Run("mark clip as lossless", func(t *testing.T) {
		timeline := &Timeline{Nodes: [][]TimelineNode{{
			CreateNode(video.NODE_VIDEO, "1", "Node1", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "2", "Node2", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "3", "Node3", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "1", "Node4", 4.2, 6.9),
			CreateNode(video.NODE_VIDEO, "1", "Node5", 4.2, 6.9)}}}

		err := timeline.ToggleLossless(0, 1)
		if err != nil {
			t.Fatal(err)
		}

		node, ok := timeline.Nodes[0][1].(*video.VideoNode)
		if !ok {
			t.Errorf("expected video node type, type: %s", node.Type())
		}
		if node.LosslessExport != true {
			t.Errorf("video clip was not marked as lossless")
		}
	})
}
