package tail

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

const (
	TransctionExt = ".transaction"
)

const (
	DateRotate = iota
	MoveRotate
)

type Tail struct {
	PathFmt  string
	Date     time.Time
	PosDir   string
	FileType int
}

type DatePos struct {
	PathFmt string
	Date    string
	FilePos int64
	Old     OutPutBuf
}
type OutPutBuf struct {
	Date    string
	Path    string
	FilePos int64
}

func (d *DatePos) GetDate() (time.Time, error) {
	return time.Parse("2006-01-02", d.Date)
}

func (d *DatePos) SetDate(t time.Time) {
	d.Date = t.Format("2006-01-02")
}

func (t *Tail) DatePos(filePathPattern string) (DatePos, error) {
	pos := DatePos{}
	posFile := path.Join(t.PosDir, path2name(filePathPattern))

	fi, err := os.Open(posFile)
	if err != nil {
		return pos, fmt.Errorf("Open: %s %s", err, posFile)
	}
	d := json.NewDecoder(fi)
	err = d.Decode(&pos)
	return pos, err
}

type PosTransaction struct {
	file string
}

// Transaction
func (t *Tail) SetDatePos(filePathPattern string, pos DatePos) (PosTransaction, error) {
	posFile := path.Join(t.PosDir, path2name(filePathPattern))

	fi, err := os.OpenFile(posFile+TransctionExt, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return PosTransaction{}, fmt.Errorf("Open:%s %s", err, posFile+TransctionExt)
	}
	e := json.NewEncoder(fi)
	err = e.Encode(&pos)
	return PosTransaction{posFile}, err
}
func (p PosTransaction) Drop() error {
	f := p.file + TransctionExt
	err := os.Remove(f)
	if err != nil {
		return fmt.Errorf("Remove transaction file error: %s %s", err, f)
	}
	return err
}
func (p PosTransaction) Commit() error {
	err := os.Rename(p.file+TransctionExt, p.file)
	if err != nil {
		return fmt.Errorf("Move transaction file error: %s %s", err, p.file)
	}
	return err
}

func path2name(p string) string {
	const escapes = "/\\?*:|\"<>[]% "
	for _, c := range escapes {
		p = strings.Replace(p, string(c), "_", -1)
	}
	return p
}

func Date2Path(p string, date time.Time) string {
	p = strings.Replace(p, "%Y", fmt.Sprintf("%04d", date.Year()), -1)
	p = strings.Replace(p, "%y", fmt.Sprintf("%02d", date.Year()%100), -1)
	p = strings.Replace(p, "%m", fmt.Sprintf("%02d", date.Month()), -1)
	p = strings.Replace(p, "%d", fmt.Sprintf("%02d", date.Day()), -1)
	return p
}
