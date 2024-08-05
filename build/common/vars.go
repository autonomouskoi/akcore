package build

import (
	"os"
	"path/filepath"
)

var (
	BaseDir string
)

func init() {
	var err error
	BaseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	BaseDir, err = filepath.Abs(filepath.Dir(BaseDir))
	if err != nil {
		panic(err)
	}
}
