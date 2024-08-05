package build

import "github.com/autonomouskoi/mageutil"

func HasCmdProtoc() error {
	return mageutil.HasExec("protoc")
}
