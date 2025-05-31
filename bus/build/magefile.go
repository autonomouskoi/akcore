package build

import (
	"path/filepath"

	build "github.com/autonomouskoi/akcore/build/common"
	"github.com/autonomouskoi/mageutil"
)

var busDir string

func init() {
	busDir = filepath.Join(build.BaseDir, "bus")
}

func Protos() error {
	return mageutil.GoProtosInDir(busDir, busDir, "module=github.com/autonomouskoi/akcore/bus")
}
