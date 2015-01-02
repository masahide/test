package tail

import (
	"os"
	"testing"
	"time"
)

func TestSetPos(t *testing.T) {
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
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	p, err := tail.DatePos()
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	if p != pos {
		t.Errorf("pos: %v\np: %v", pos, p)
	}
	/*
		err = os.Remove(path2name(posfile))
		if err != nil {
			t.Errorf("got %v\nwant %v", err, nil)
		}
	*/
}

func TestOverWriteSetDatePos(t *testing.T) {
	tail := &Tail{
		PosDir:  "./",
		PathFmt: "/fuga/%Y%m%d/hoge/file.log",
	}

	pos := DatePos{
		PathFmt: "hoge",
		FilePos: 1,
	}
	pos.SetDate(time.Now())
	tail.SetDatePos(pos)
	tail.SetDatePos(pos)
	pt, _ := tail.SetDatePos(pos)
	pt.Commit()
	p, _ := tail.DatePos()
	if p != pos {
		t.Errorf("got %v\nwant %v", pos, p)
	}
	err := os.Remove(path2name(tail.PathFmt))
	if err != nil {
		t.Errorf("got %v\nwant %v", err, nil)
	}
}

type date2pathTestpair struct {
	date string
	p    string
	want string
}

var date2pathTests = []date2pathTestpair{
	{"2013-06-19 21:54:23 +0900", "/hoge%Y%m%d/_%Y%m%dhoge", "/hoge20130619/_20130619hoge"},
	{"2013-06-19 21:54:23 +0900", "/hoge%y%m%d/_%Y%m%dhoge", "/hoge130619/_20130619hoge"},
	{"1993-06-19 21:54:23 +0900", "/hoge%y%m%d/_%Y%m%dhoge", "/hoge930619/_19930619hoge"},
}

func TestDate2path(t *testing.T) {
	for _, pair := range date2pathTests {
		d, err := time.Parse("2006-01-02 15:04:05 -0700", pair.date)
		if err != nil {
			t.Error(err)
		}
		v := Date2Path(pair.p, d)
		if v != pair.want {
			t.Error(
				"For:", pair.date,
				"want:", pair.want,
				"got:", v,
			)
		}
	}
}
