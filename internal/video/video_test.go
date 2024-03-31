package video

import (
	"math"
	"testing"
)

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

func TestDeleteRIDByReferences(t *testing.T) {
	t.Run("delete references of a rid", func(t *testing.T) {
		expected := &Timeline{
			VideoNodes: []VideoNode{createVideoNode("2", 4.2, 6.9),
				createVideoNode("3", 4.2, 6.9),
			},
		}

		tl := NewTimeline()
		_, _ = tl.Insert("1", 4.2, 6.9, 0)
		_, _ = tl.Insert("2", 4.2, 6.9, 1)
		_, _ = tl.Insert("3", 4.2, 6.9, 2)
		_, _ = tl.Insert("1", 4.2, 6.9, 3)
		_, _ = tl.Insert("1", 4.2, 6.9, 4)

		err := tl.DeleteRIDReferences("1")
		if err != nil {
			t.Errorf("failed to perform reference deletion")
		}

		if len(tl.VideoNodes) != len(expected.VideoNodes) {
			t.Errorf("the timeline lenghts do not match")
		}

		for i := range tl.VideoNodes {
			if tl.VideoNodes[i].RID != expected.VideoNodes[i].RID &&
				math.Abs(tl.VideoNodes[i].Start-expected.VideoNodes[i].Start) <= 0.01 &&
				math.Abs(tl.VideoNodes[i].End-expected.VideoNodes[i].End) <= 0.01 {
				t.Errorf("the timelines do not have the same order")
			}
		}

	})
}
