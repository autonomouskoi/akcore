package build

import (
	"github.com/autonomouskoi/mageutil"
	"github.com/magefile/mage/mg"
)

func GoProto(dest, src string) error {
	mg.Deps(HasCmdProtoc)
	return mageutil.GoProto(dest, src, BaseDir, "module=github.com/autonomouskoi/akcore")
}
