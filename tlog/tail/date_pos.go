package tail

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

type DatePos struct {
	PathFmt string
	Path    string
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

// Transaction
type PosTransaction struct {
	file string
}

func (t *DateReader) PosTransaction() (PosTransaction, error) {
	posFile := path.Join(t.PosDir, path2name(t.PathFmt))

	tFileName := posFile + TransctionExt
	fi, err := os.OpenFile(tFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return PosTransaction{}, fmt.Errorf("Open:%s %s", err, tFileName)
	}
	pos, err := t.DatePos()
	if err != nil {
		return PosTransaction{}, err
	}
	pos.Old.Date = pos.Date
	pos.Old.FilePos = pos.FilePos
	pos.Old.Path = pos.Path
	e := json.NewEncoder(fi)
	err = e.Encode(&pos)
	fi.Close()
	return PosTransaction{posFile}, err
}
func (t *DateReader) SetDatePos(pos DatePos) (PosTransaction, error) {
	posFile := path.Join(t.PosDir, path2name(t.PathFmt))

	fi, err := os.OpenFile(posFile+TransctionExt, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return PosTransaction{}, fmt.Errorf("Open:%s %s", err, posFile+TransctionExt)
	}
	e := json.NewEncoder(fi)
	err = e.Encode(&pos)
	fi.Close()
	return PosTransaction{posFile}, err
}
func (p PosTransaction) Commit() error {
	err := os.Rename(p.file+TransctionExt, p.file)
	if err != nil {
		return fmt.Errorf("Move transaction file error: %s %s", err, p.file)
	}
	return err
}

func (p PosTransaction) Drop() error {
	f := p.file + TransctionExt
	err := os.Remove(f)
	if err != nil {
		return fmt.Errorf("Remove transaction file error: %s %s", err, f)
	}
	return err
}
