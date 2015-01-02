package tail

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	TransctionExt = ".transaction"

	DateRotate = iota
	MoveRotate
)

// pathFmtのフォーマット文字列をアンダースコアに置換
func path2name(p string) string {
	const escapes = "/\\?*:|\"<>[]% "
	for _, c := range escapes {
		p = strings.Replace(p, string(c), "_", -1)
	}
	return p
}

// Pathの日付フォーマットに日付を適用
func Date2Path(p string, date time.Time) string {
	p = strings.Replace(p, "%Y", fmt.Sprintf("%04d", date.Year()), -1)
	p = strings.Replace(p, "%y", fmt.Sprintf("%02d", date.Year()%100), -1)
	p = strings.Replace(p, "%m", fmt.Sprintf("%02d", date.Month()), -1)
	p = strings.Replace(p, "%d", fmt.Sprintf("%02d", date.Day()), -1)
	return p
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
