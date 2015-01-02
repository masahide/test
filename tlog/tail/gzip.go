package tail

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}
type GzReadSeeker struct {
	file    *os.File
	decodeF *gzip.Reader
}

func (grs *GzReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if whence != 1 {
		return 0, errors.New("Seek err: Not support whence.")
	}
	return io.CopyN(ioutil.Discard, grs, offset)
}

func (grs *GzReadSeeker) Close() error {
	errd := grs.decodeF.Close()
	errf := grs.file.Close()
	if errd != nil {
		return errd
	}
	if errf != nil {
		return errf
	}
	return nil
}

func (grs *GzReadSeeker) Read(p []byte) (n int, err error) {
	return grs.decodeF.Read(p)
}

func GzOpen(file string) (*GzReadSeeker, error) {
	var grs GzReadSeeker
	var err error
	grs.file, err = os.OpenFile(file, os.O_RDONLY, 0600)
	if err != nil {
		return &grs, err
	}
	grs.decodeF, err = gzip.NewReader(grs.file)
	if err != nil {
		grs.file.Close()
		return &grs, err
	}
	return &grs, nil
}
