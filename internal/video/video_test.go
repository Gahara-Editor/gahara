package video

import (
	"math"
	"testing"
)

func mockTl() *Timeline {
	return &Timeline{VideoNodes: []VideoNode{
		createVideoNode("1", "Node", 4.2, 6.9),
		createVideoNode("2", "Node", 4.2, 6.9),
		createVideoNode("3", "Node", 4.2, 6.9),
		createVideoNode("1", "Node", 4.2, 6.9),
		createVideoNode("1", "Node", 4.2, 6.9)}}

}

func TestGetFilename(t *testing.T) {
	t.Run("get filename", func(t *testing.T) {
		got := GetFilename("../dir/dir/path/filename.txt")
		expected := "filename.txt"
		if got != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})

	t.Run("empty path should return empty filename", func(t *testing.T) {
		got := GetFilename("")
		expected := ""
		if got != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})
}

func TestGetNameAndExtension(t *testing.T) {
	t.Run("get name and extension", func(t *testing.T) {
		name, ext, err := GetNameAndExtension("myfile.mp4")
		expectedName, expectedExt := "myfile", "mp4"

		if name != expectedName || ext != expectedExt {
			t.Errorf("got {name: %q, ext :%q}, expected {name: %q, ext :%q}", name, ext, expectedName, expectedExt)
		}
		if err != nil {
			t.Errorf("error getting name and extension %s", err.Error())
		}
	})
}

func TestTimelineInsert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		element  VideoNode
		idx      int
		tl       *Timeline
		expected *Timeline
	}{
		{
			name:    "Insert an element in the middle",
			element: createVideoNode("middle", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     2,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("middle", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name:    "Insert an element at the beginning",
			element: createVideoNode("first", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     0,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("first", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name: "Insert an element at the end",
			tl:   mockTl(),
			idx:  5,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("last", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name:    "Index out out bounds (slice size = 4, idx = -1)",
			element: createVideoNode("fail", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     -1,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name:    "Index out out bounds (slice size = 4, idx = 10)",
			element: createVideoNode("fail", "Node", 4.2, 6.9),
			tl:      mockTl(),
			idx:     6,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.tl.Insert(tt.element.RID, tt.element.Name, tt.element.Start, tt.element.End, tt.idx)
		if tt.idx < 0 || tt.idx > len(tt.tl.VideoNodes) {
			if err == nil {
				t.Errorf("index %d is out of bounds and it should have failed", tt.idx)
			}
			return
		}

		if err != nil {
			t.Error("could not perform delete operation: ", err.Error())
		}

		if len(tt.tl.VideoNodes) != len(tt.expected.VideoNodes) {
			t.Errorf("the timeline lenghts do not match: (expected: %d, got: %d)", len(tt.expected.VideoNodes), len(tt.tl.VideoNodes))
		}

		for i := range tt.tl.VideoNodes {
			if tt.tl.VideoNodes[i].RID != tt.expected.VideoNodes[i].RID &&
				tt.tl.VideoNodes[i].Name != tt.expected.VideoNodes[i].Name &&
				math.Abs(tt.tl.VideoNodes[i].Start-tt.expected.VideoNodes[i].Start) <= 0.01 &&
				math.Abs(tt.tl.VideoNodes[i].End-tt.expected.VideoNodes[i].End) <= 0.01 {
				t.Errorf("the timelines do not have the same order")
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
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name: "Delete the first element",
			tl:   mockTl(),
			idx:  0,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name: "Delete the last element",
			tl:   mockTl(),
			idx:  0,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name: "Index out out bounds (slice size = 4, idx = -1)",
			tl:   mockTl(),
			idx:  -1,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
		{
			name: "Index out out bounds (slice size = 4, idx = 10)",
			tl:   mockTl(),
			idx:  10,
			expected: &Timeline{VideoNodes: []VideoNode{
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
				createVideoNode("1", "Node", 4.2, 6.9),
			},
			},
		},
	}

	for _, tt := range tests {
		err := tt.tl.Delete(tt.idx)
		if tt.idx < 0 || tt.idx > len(tt.tl.VideoNodes) {
			if err == nil {
				t.Errorf("index %d is out of bounds and it should have failed", tt.idx)
			}
			return
		}

		if err != nil {
			t.Error("could not perform delete operation: ", err.Error())
		}

		if len(tt.tl.VideoNodes) != len(tt.expected.VideoNodes) {
			t.Errorf("the timeline lenghts do not match: (expected: %d, got: %d)", len(tt.expected.VideoNodes), len(tt.tl.VideoNodes))
		}

		for i := range tt.tl.VideoNodes {
			if tt.tl.VideoNodes[i].RID != tt.expected.VideoNodes[i].RID &&
				tt.tl.VideoNodes[i].Name != tt.expected.VideoNodes[i].Name &&
				math.Abs(tt.tl.VideoNodes[i].Start-tt.expected.VideoNodes[i].Start) <= 0.01 &&
				math.Abs(tt.tl.VideoNodes[i].End-tt.expected.VideoNodes[i].End) <= 0.01 {
				t.Errorf("the timelines do not have the same order")
			}
		}

	}
}

func TestDeleteRIDByReferences(t *testing.T) {
	t.Run("delete references of a rid", func(t *testing.T) {
		expected := &Timeline{
			VideoNodes: []VideoNode{createVideoNode("2", "Node", 4.2, 6.9),
				createVideoNode("3", "Node", 4.2, 6.9),
			},
		}

		tl := &Timeline{VideoNodes: []VideoNode{
			createVideoNode("1", "Node", 4.2, 6.9),
			createVideoNode("2", "Node", 4.2, 6.9),
			createVideoNode("3", "Node", 4.2, 6.9),
			createVideoNode("1", "Node", 4.2, 6.9),
			createVideoNode("1", "Node", 4.2, 6.9)}}

		err := tl.DeleteRIDReferences("1")
		if err != nil {
			t.Errorf("failed to perform reference deletion")
		}

		if len(tl.VideoNodes) != len(expected.VideoNodes) {
			t.Errorf("the timeline lenghts do not match: (expected: %d, got: %d)", len(expected.VideoNodes), len(tl.VideoNodes))
		}

		for i := range tl.VideoNodes {
			if tl.VideoNodes[i].RID != expected.VideoNodes[i].RID &&
				tl.VideoNodes[i].Name != expected.VideoNodes[i].Name &&
				math.Abs(tl.VideoNodes[i].Start-expected.VideoNodes[i].Start) <= 0.01 &&
				math.Abs(tl.VideoNodes[i].End-expected.VideoNodes[i].End) <= 0.01 {
				t.Errorf("the timelines do not have the same order")
			}
		}

	})
}

func TestMarkLossless(t *testing.T) {
	t.Run("mark clip as lossless", func(t *testing.T) {
		timeline := &Timeline{VideoNodes: []VideoNode{
			createVideoNode("1", "Node1", 4.2, 6.9),
			createVideoNode("2", "Node2", 4.2, 6.9),
			createVideoNode("3", "Node3", 4.2, 6.9),
			createVideoNode("1", "Node4", 4.2, 6.9),
			createVideoNode("1", "Node5", 4.2, 6.9)}}

		err := timeline.ToggleLossless(1)
		if err != nil {
			t.Fatal(err)
		}

		if timeline.VideoNodes[1].LosslessExport != true {
			t.Errorf("video clip was not marked as lossless")
		}
	})
}
