package core

import "os"

var (
	KlineFilePath  string = "data/"
	ReportFilePath string = "reports/"
	ChartFilePath  string = "charts/"
)

func init() {
	createDirectory(KlineFilePath)
	createDirectory(ReportFilePath)
	createDirectory(ChartFilePath)
}

func createDirectory(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)

		if err != nil {
			panic(err)
		}
	}
}
