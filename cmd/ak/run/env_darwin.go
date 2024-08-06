package run

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed icon.png
var IconBytes []byte

func AppPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "AutonomousKoi"), nil
}

func ShowFolder(path string) error {
	return exec.Command("open", path).Run()
}
