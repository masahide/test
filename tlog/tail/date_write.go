package tail

import (
	"os"
	"path/filepath"
	"time"
)

type DateWriter struct {
	PathFmt  string
	OutDir   string
	Interval time.Duration
	Daemon   string
	Type     string
	Name     string

	dateTime time.Time
	file     *os.File
}

//        saveTime := time.Unix(line.Time, 0).Truncate(s.RotateTime)

func (t *DateWriter) mkFilePath() string {
	return filepath.Join(t.OutDir, t.Daemon+"."+t.Type+","+t.Name+"."+t.dateTime.Format("200601021504"))
}

func (t *DateWriter) OpenCloseDate(dt time.Time) error {
	var err error
	dt = dt.Truncate(t.Interval)
	if t.file != nil && !dt.Equal(t.dateTime) {
		err = t.file.Close()
		if err != nil {
			return err
		}
		t.file, err = os.OpenFile(t.mkFilePath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	return err
}

func (t *DateWriter) Date() time.Time {
	return t.dateTime
}
func (t *DateWriter) SetDate(date time.Time) {
	t.dateTime = date
}

func (t *DateWriter) Write(p []byte) (n int, err error) {
	return t.file.Write(p)
}
