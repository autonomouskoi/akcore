package internal

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"

	"github.com/autonomouskoi/mageutil"
)

var (
	baseDir string
)

func init() {
	var err error
	baseDir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	baseDir = filepath.Join(baseDir, "..", "internal")
}

func All() {
	mg.Deps(
		Dev,
	)
}

func Dev() {
	mg.Deps(
		GoProtos,
	)
}

func GoProtos() error {
	return mageutil.GoProtosInDir(baseDir, "module=github.com/autonomouskoi/akcore/internal")
}
