package app

import (
	"path"
	"time"
)

type Xlsx struct {
	XlsxDir string
}

// XlsxFilePath
// name with xlsx file suffix
func (x *Xlsx) XlsxFilePath(fileName string) string {
	return path.Join(x.XlsxDir, fileName)
}

func (x *Xlsx) XlsxFilePathWithName(name string) string {
	return x.XlsxFilePath(name + ".xlsx")
}

func (x *Xlsx) XlsxFilePathWithNameWithTime(name string) string {
	return x.XlsxFilePathWithName(name + time.Now().Format("20060102_1504"))
}

func (x *Xlsx) XlsxQueryWithNameWithTime(name string) string {
	return x.XlsxFilePathWithNameWithTime("_q_" + name)
}
