package tail

import (
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	DateReader
	tail := &Tail{
		PosDir:  "./",
		PathFmt: "/fuga/%Y%m%d/hoge/file.log",
	}

	pos := DatePos{
		PathFmt: "hoge",
		FilePos: 1,
	}
	pos.SetDate(time.Now())

	pt, err := tail.SetDatePos(pos)
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	err = pt.Commit()
	if err != nil {
	}
}
