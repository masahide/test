package tail

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/cenkalti/backoff"
)

const (
	delay = 30 * time.Second
)

type DateReader struct {
	PathFmt   string
	PosDir    string
	StartTail bool

	date     time.Time
	file     ReadSeekCloser
	lastByte byte
	retry    *backoff.ExponentialBackOff
}

func (t *DateReader) Date() time.Time {
	return t.date
}

func (t *DateReader) SetDate(date time.Time) {
	t.date = date
}

func (t *DateReader) initRetry() {
	t.retry = &backoff.ExponentialBackOff{
		InitialInterval:     10 * time.Microsecond,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         1 * time.Second,
		MaxElapsedTime:      5 * time.Second,
		Clock:               backoff.SystemClock,
	}
}

func (t *DateReader) Open() error {
	if t.file != nil {
		return nil
	}
	t.initRetry()
	pos, err := t.DatePos()
	if err != nil {
		return err
	}
	if pos == (DatePos{}) { // posファイルがない場合は本日分のログファイルを指定
		if t.StartTail {
			return t.OpenTail()
		}
		return t.OpenDate(time.Now(), 0)
	}
	t.date, err = pos.GetDate()
	if err != nil {
		return err
	}
	return t.openFile(pos)
}

// 日付のファイルを開く
func (t *DateReader) OpenTail() error {
	err := t.OpenDate(time.Now(), 0)
	if err != nil {
		return err
	}
	_, err = t.file.Seek(0, os.SEEK_END)
	return err
}

// 日付のファイルを開く
func (t *DateReader) OpenDate(date time.Time, FilePos int64) error {
	t.initRetry()
	pos := DatePos{
		PathFmt: t.PathFmt,
		FilePos: FilePos,
	}
	pos.SetDate(date)

	return t.openFile(pos)
}

// ポジション情報からファイルを開く
func (t *DateReader) openFile(pos DatePos) error {
	var err error
	t.date, err = pos.GetDate()
	if err != nil {
		return err
	}
	filePath := Date2Path(pos.PathFmt, t.date)
	if Exists(filePath) {
		t.file, err = os.OpenFile(filePath, os.O_RDONLY, 0600)
	} else {
		t.file, err = GzOpen(filePath + ".gz") // 生ファイルがなければgzファイルを開く
	}
	if err != nil {
		return err
	}
	if pos.FilePos != 0 {
		_, err = t.file.Seek(pos.FilePos, os.SEEK_CUR)
	}
	return err
}

func (t *DateReader) FileTime() time.Time {
	today := time.Now().Truncate(24 * time.Hour)
	if t.date.Equal(today) { // 当日
		return time.Now()
	}
	if t.date.Equal(today.Add(-24 * time.Hour)) { // 昨日
		if time.Since(today) <= delay { // 日替わり直後
			return t.date.Add(24*time.Hour - time.Second)
		}
	}
	return t.date
}

// 開いているファイルが最新のログファイルかどうか
func (t *DateReader) IsRealtimeRead() bool {
	today := time.Now().Truncate(24 * time.Hour)
	if t.date.Equal(today) { // 当日
		return true
	}
	if t.date.Equal(today.Add(-24 * time.Hour)) { // 昨日
		return time.Since(today) <= delay // 日替わり直後
	}
	return false
}

// io.ReaderのReadを満たす
func (t *DateReader) Read(p []byte) (int, error) {
	nw := 0
	err := t.Open()
	if err != nil {
		return nw, err
	}
	if !t.IsRealtimeRead() {
		return t.file.Read(p) // 十分に過去
	}
	err = backoff.Retry(func() error {
		var bl int
		bl, err = t.file.Read(p)
		nw += bl
		if bl > 0 {
			t.lastByte = p[bl-1]
			if len(p) == bl || t.lastByte == byte('\n') {
				return nil
			}
		} else if t.lastByte == byte('\n') {
			return nil
		}
		p = p[bl:]
		return errors.New("Did not find a new line code.")
	}, t.retry)
	return nw, err
}

func (t *DateReader) DatePos() (DatePos, error) {
	pos := DatePos{}
	posFile := path.Join(t.PosDir, path2name(t.PathFmt))

	if !Exists(posFile) {
		return pos, nil // ファイルがない場合は空を返す
	}

	fi, err := os.Open(posFile)
	if err != nil {
		return pos, fmt.Errorf("Open: %s %s", err, posFile)
	}
	d := json.NewDecoder(fi)
	err = d.Decode(&pos)
	fi.Close()
	return pos, err
}
