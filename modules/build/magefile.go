package modules

import (
	"fmt"
	"path/filepath"

	"github.com/magefile/mage/mg"

	build "github.com/autonomouskoi/akcore/build/common"
	bus "github.com/autonomouskoi/akcore/bus/build"
)

var modulesDir string

func init() {
	modulesDir = filepath.Join(build.BaseDir, "modules")
}

func Main() {
	mg.Deps(
		bus.Protos,
		Protos,
	)
}

func Protos() error {
	for _, baseName := range []string{
		"config",
		"control",
		"manifest",
	} {
		src := filepath.Join(modulesDir, baseName+".proto")
		dest := filepath.Join(modulesDir, baseName+".pb.go")
		if err := build.GoProto(dest, src); err != nil {
			return fmt.Errorf("building %s: %w", dest, err)
		}
	}
	return nil
}
