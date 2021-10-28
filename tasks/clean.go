package tasks

import (
	"os"
	"statika/util"
)

func Clean(outputDir string) {
	err := os.RemoveAll(outputDir)
	util.Check(err)

	err = os.MkdirAll(outputDir, 0777)
	util.Check(err)
}
