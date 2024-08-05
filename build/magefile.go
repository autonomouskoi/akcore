//go:build mage
// +build mage

package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	build "github.com/autonomouskoi/akcore/build/common"
	modules "github.com/autonomouskoi/akcore/modules/build"
	web "github.com/autonomouskoi/akcore/web/build"
	"github.com/autonomouskoi/mageutil"
)

func Clean() {
	mg.Deps(web.Clean)
}

func Dev() {
	mg.Deps(
		Modules,
		WebContent,
	)
}

func DevRun() error {
	mg.Deps(
		Dev,
	)
	cmdPath := filepath.Join(build.BaseDir, "cmd", "aktrackstar")
	return sh.RunWith(
		map[string]string{
			"AK_WEB_CONTENT":                 filepath.Join(build.BaseDir, "web", "content", "out"),
			"AK_TRACKSTAR_CONTENT":           filepath.Join(build.BaseDir, "..", "trackstar", "web"),
			"AK_TRACKSTAR_OVERLAY_CONTENT":   filepath.Join(build.BaseDir, "..", "trackstar", "overlay", "web"),
			"AK_TRACKSTAR_STAGELINQ_CONTENT": filepath.Join(build.BaseDir, "..", "trackstar", "stagelinq", "web"),
		},
		"go", "run", cmdPath,
	)
}

func Modules() {
	mg.Deps(modules.Main)
}

func WebContent() {
	mg.Deps(web.Content)
}

func WebZip() {
	mg.Deps(web.Zip)
}

var releaseTimestamp string
var distDir string
var mainPath string

func Release() {
	mg.SerialDeps(
		ReleaseMac,
		ReleaseWin,
	)
}

func ReleaseDeps() error {
	releaseTimestamp = time.Now().Format("20060102-1504")
	distDir = filepath.Join(build.BaseDir, "dist")
	if err := sh.Rm(distDir); err != nil {
		return fmt.Errorf("removing %s: %w", distDir, err)
	}
	if err := mageutil.Mkdir(distDir); err != nil {
		return fmt.Errorf("creating %s: %w", distDir, err)
	}
	mainPath = filepath.Join(build.BaseDir, "cmd", "aktrackstar")
	mg.Deps(
		Modules,
		WebZip,
	)
	return nil
}

func ReleaseWin() error {
	mg.Deps(ReleaseDeps)
	baseName := "aktrackstar-win-" + releaseTimestamp
	outPath := filepath.Join(distDir, baseName+".exe")
	err := sh.RunWith(map[string]string{
		"CGO_ENABLED": "1",
		"CGO_CFLAGS":  "-I/mingw64/include",
		"MSYSTEM":     "MINGW64",
	},
		"go", "build",
		"-o", outPath,
		"-ldflags", "-H=windowsgui",
		mainPath,
	)
	if err != nil {
		return fmt.Errorf("building %s: %w", outPath, err)
	}
	zipPath := filepath.Join(distDir, baseName+".zip")
	err = mageutil.ZipFiles(zipPath, map[string]string{
		filepath.Join(build.BaseDir, "LICENSE"): "LICENSE",
		outPath:                                 baseName + ".exe",
	})
	return err
}

func ReleaseMac() error {
	mg.Deps(ReleaseDeps)
	baseName := "aktrackstar-mac-" + releaseTimestamp
	outPath := filepath.Join(distDir, baseName)
	err := sh.RunWith(map[string]string{},
		"go", "build",
		"-o", outPath,
		"-ldflags", "-s -w",
		mainPath,
	)
	if err != nil {
		return fmt.Errorf("building %s: %w", outPath, err)
	}

	dmgTmplPath := filepath.Join(mainPath, "run", "AK-tmpl.dmg.gz")
	tmplFH, err := os.Open(dmgTmplPath)
	if err != nil {
		return fmt.Errorf("opening DMG template %s: %w", dmgTmplPath, err)
	}
	defer tmplFH.Close()
	gzR, err := gzip.NewReader(tmplFH)
	if err != nil {
		return fmt.Errorf("creating DMG template decompressor: %w", err)
	}
	dmgFilePath := filepath.Join(distDir, "AK-tmpl.dmg")
	dmgFH, err := os.Create(dmgFilePath)
	if err != nil {
		return fmt.Errorf("creating DMG file %s: %w", dmgFilePath, err)
	}
	defer dmgFH.Close()
	if _, err := io.Copy(dmgFH, gzR); err != nil {
		return fmt.Errorf("decompressing DMG file: %w", err)
	}
	if err := dmgFH.Sync(); err != nil {
		return fmt.Errorf("syncing DMG file: %w", err)
	}

	// resize to hold the executable + 1MB overhead
	stat, err := os.Stat(outPath)
	if err != nil {
		return fmt.Errorf("statting executable: %w", err)
	}
	err = sh.Run("hdiutil", "resize", "-size", strconv.Itoa(int(stat.Size())+(1024*1024)), dmgFilePath)
	if err != nil {
		return fmt.Errorf("resizing DMG file: %w", err)
	}

	// mount the dmg
	appDir := filepath.Join(distDir, "mac")
	if err := mageutil.Mkdir(appDir); err != nil {
		return fmt.Errorf("creating app dir: %w", err)
	}
	err = sh.Run("hdiutil", "attach", dmgFilePath, "-noautoopen", "-mountpoint", appDir)
	if err != nil {
		return fmt.Errorf("attaching DMG: %w", err)
	}
	detached := false
	defer func() {
		if !detached {
			sh.Run("hdiutil", "detach", appDir+"/")
		}
	}()

	// copy stuff
	appExecPath := filepath.Join(appDir, "AutonomousKoi.app", "Contents", "MacOS", "ak")
	if err := sh.Copy(appExecPath, outPath); err != nil {
		return fmt.Errorf("copying app executable: %w", err)
	}
	if err := os.Chmod(appExecPath, 0555); err != nil {
		return fmt.Errorf("setting app executable permissions: %w", err)
	}
	licDestPath := filepath.Join(appDir, "LICENSE")
	licSrcPath := filepath.Join(build.BaseDir, "LICENSE")
	if err := sh.Copy(licDestPath, licSrcPath); err != nil {
		return fmt.Errorf("copying LICENSE: %w", err)
	}

	// detach, compress
	if err := sh.Run("hdiutil", "detach", appDir+"/"); err != nil {
		return fmt.Errorf("detaching DMG %s: %w", appDir, err)
	}
	detached = true
	err = sh.Run("hdiutil", "convert", dmgFilePath,
		"-format", "UDZO",
		"-imagekey", "zlib-level=9",
		"-o", filepath.Join(distDir, "AutonomousKoi-"+releaseTimestamp+".dmg"),
	)
	if err != nil {
		return fmt.Errorf("compressing DMG: %w", err)
	}

	return nil
}