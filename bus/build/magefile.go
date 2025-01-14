package build

import (
	"fmt"
	"path/filepath"

	build "github.com/autonomouskoi/akcore/build/common"
)

var busDir string

func init() {
	busDir = filepath.Join(build.BaseDir, "bus")
}

func Protos() error {
	for _, baseName := range []string{
		"bus",
		"direct",
	} {
		src := filepath.Join(busDir, baseName+".proto")
		dest := filepath.Join(busDir, baseName+".pb.go")
		if err := build.GoProto(dest, src); err != nil {
			return fmt.Errorf("building %s: %w", dest, err)
		}
	}
	return nil
}
