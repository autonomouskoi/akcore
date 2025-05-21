package modules

import (
	"path/filepath"

	"github.com/magefile/mage/mg"

	build "github.com/autonomouskoi/akcore/build/common"
	bus "github.com/autonomouskoi/akcore/bus/build"
	"github.com/autonomouskoi/mageutil"
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
	return mageutil.GoProtosInDir(modulesDir, modulesDir, "")
}
