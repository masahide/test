package tail

import (
	"fmt"
	"io"
	"time"
)

type DateConfig struct {
	PathFmt   string
	PosDir    string
	OutDir    string
	Daemon    string
	Type      string
	Name      string
	Interval  int // time.Second
	StartDt   string
	StartTail bool
}

type DateTail struct {
	In  DateReader
	Out DateWriter
	Log *Logger
}

func (t *DateTail) TailF() {
	for {
		t.Tail()

	}
}

func (t *DateTail) Tail() {
	t.In.Open()
	ft := t.In.FileTime()
	t.Out.OpenCloseDate(ft)
	io.Copy(&t.Out, &t.In)
}

func NewDateTail(c DateConfig) (*DateTail, error) {

	dt := DateTail{
		In: DateReader{
			PathFmt:   c.PathFmt,
			PosDir:    c.PosDir,
			StartTail: c.StartTail,
		},
		Out: DateWriter{
			PathFmt:  c.PathFmt,
			OutDir:   c.OutDir,
			Interval: time.Duration(c.Interval) * time.Second,
			Daemon:   c.Daemon,
			Type:     c.Type,
			Name:     c.Name,
		},
		Log: NewLogger(),
	}

	pos, err := dt.In.DatePos()
	if err != nil {
		return nil, err
	}
	if pos == (DatePos{}) { // posファイルがない場合
		var err error
		date := time.Now().Truncate(24 * time.Hour) //とりあえずtodayを入れる
		if c.StartDt != "" {                        // start dateの指定がある場合はその日付を入れる
			date, err = time.Parse(
				"2006-01-02 15:04:05", //"2006-01-02 15:04:05 -0700",    // scan format
				c.StartDt,
			)
			if err != nil {
				return &dt, fmt.Errorf("Config.StartDt time.Parse err: %s", err)
			}
			date = date.Truncate(24 * time.Hour)
		}
		dt.In.SetDate(date)
	}

	return &dt, nil

}
