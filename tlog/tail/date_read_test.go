package tail

import "testing"

func TestOpen(t *testing.T) {
	dr := DateReader{
		PathFmt:   "./date_read_test.go",
		PosDir:    "./",
		StartTail: true,
	}

	err := dr.Open()
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
}
