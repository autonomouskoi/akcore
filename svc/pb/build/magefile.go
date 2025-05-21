package svc

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

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
	baseDir = filepath.Join(baseDir, "..", "svc", "pb")
}

func All() {
	mg.Deps(
		Dev,
	)
}

func Clean() error {
	matches, err := fs.Glob(os.DirFS(baseDir), "*.pb.go")
	if err != nil {
		return fmt.Errorf("globbing %s: %w", baseDir, err)
	}
	for _, match := range matches {
		if err := sh.Rm(match); err != nil {
			return fmt.Errorf("removing %s: %w", match, err)
		}
	}
	return nil
}

func Dev() {
	mg.Deps(
		GoProtos,
	)
}

func GoProtos() error {
	return mageutil.GoProtosInDir(baseDir, baseDir, "module=github.com/autonomouskoi/akcore/svc/pb/svc")
}
