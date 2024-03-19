package utils

import "testing"

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
