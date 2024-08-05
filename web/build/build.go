package web

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"

	build "github.com/autonomouskoi/akcore/build/common"
	bus "github.com/autonomouskoi/akcore/bus/build"
	"github.com/autonomouskoi/mageutil"
)

var (
	webContentDir    string
	webContentOutDir string
	webContentZip    string
)

func init() {
	webContentDir = filepath.Join(build.BaseDir, "web", "content")
	webContentOutDir = filepath.Join(webContentDir, "out")
	webContentZip = filepath.Join(build.BaseDir, "web", "web.zip")
}

func Clean() {
	for _, path := range []string{
		webContentOutDir,
		webContentZip,
	} {
		sh.Rm(path)
	}
}

func Content() {
	mg.Deps(
		Dirs,
		bus.Protos,
	)
	mg.Deps(
		FilesOther,
		ProtoLib,
		Protos,
		SrcCopy,
		TS,
	)
}

func Dirs() error {
	return mageutil.Mkdir(webContentOutDir)
}

func FilesOther() error {
	pixiMJS := filepath.Join(webContentDir, "node_modules/pixi.js/dist/pixi.min.mjs")
	pixiMJSMap := filepath.Join(webContentDir, "node_modules/pixi.js/dist/pixi.min.mjs.map")
	err := mageutil.HasFiles(
		pixiMJS, pixiMJSMap,
	)
	if err != nil {
		return err
	}
	mg.Deps(Dirs)
	return mageutil.CopyFiles(map[string]string{
		pixiMJS:    filepath.Join(webContentOutDir, "pixi.js"),
		pixiMJSMap: filepath.Join(webContentOutDir, "pixi.mjs.map"),
	})
}

func ProtoLib() error {
	srcDir := filepath.Join(webContentDir,
		"node_modules/@bufbuild/protobuf/dist/esm/",
	)
	if err := mageutil.HasFiles(srcDir); err != nil {
		return err
	}
	mg.Deps(Dirs)
	destDir := filepath.Join(webContentOutDir, "protobuf")
	newer, err := target.Dir(destDir, srcDir)
	if err != nil {
		return fmt.Errorf("testing %s vs %s: %w", srcDir, destDir, err)
	}
	if !newer {
		return nil
	}
	mageutil.VerboseF("copying %s -> %s\n", srcDir, destDir)
	return mageutil.CopyRecursively(destDir, srcDir)
}

func Protos() error {
	mg.Deps(build.HasCmdProtoc)
	plugin := filepath.Join(webContentDir,
		"node_modules/.bin/protoc-gen-es",
	)
	if runtime.GOOS == "windows" {
		plugin += ".cmd"
	}
	if err := mageutil.HasFiles(plugin); err != nil {
		return err
	}
	mg.Deps(Dirs)
	protoDestDir := filepath.Join(webContentOutDir, "pb")
	if err := mageutil.Mkdir(protoDestDir); err != nil {
		return fmt.Errorf("creating %s: %w", protoDestDir, err)
	}
	for _, srcFile := range []string{
		//"autonomouskoi.proto",
		"bus/bus.proto",
		"modules/control.proto",
		"modules/manifest.proto",
		//"twitch.proto",
		//"modules/magic/magic.proto",
		//"modules/twitchemotefx/twitchemotefx.proto",
	} {
		baseName := strings.TrimSuffix(filepath.Base(srcFile), ".proto")
		destFile := filepath.Join(protoDestDir, baseName+"_pb.js")
		srcFile = filepath.Join(build.BaseDir, srcFile)
		newer, err := target.Path(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("testing %s vs %s: %w", srcFile, destFile, err)
		}
		if !newer {
			continue
		}
		mageutil.VerboseF("generating proto %s -> %s\n", srcFile, destFile)
		err = sh.Run("protoc",
			"--plugin", "protoc-gen-es="+plugin,
			"-I", build.BaseDir,
			"--es_out", protoDestDir,
			srcFile,
		)
		if err != nil {
			return fmt.Errorf("generating proto %s -> %s\n: %w", srcFile, destFile, err)
		}
	}
	return nil
}

func SrcCopy() error {
	return mageutil.CopyInDir(webContentOutDir, webContentDir,
		"index.html", "ui.html",
		"OBS_Studio_Logo.svg", "links-line.svg", "equalizer-line.svg",
		"favicon.ico",
	)
}

func TS() error {
	mg.Deps(Protos)
	dirEntries, err := os.ReadDir(webContentDir)
	if err != nil {
		return fmt.Errorf("listing %s: %w", webContentDir, err)
	}
	paths := map[string]string{}
	for _, entry := range dirEntries {
		if entry.Type() == os.ModeDir {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".ts") {
			continue
		}
		baseName := strings.TrimSuffix(name, ".ts")
		destFile := filepath.Join(webContentOutDir, baseName+".js")
		srcFile := filepath.Join(webContentDir, name)
		paths[srcFile] = destFile
	}
	newer, err := mageutil.Newer(paths)
	if err != nil {
		return err
	}
	tsc := filepath.Join(webContentDir, "node_modules", ".bin", "tsc")
	if runtime.GOOS == "windows" {
		tsc += ".cmd"
	}
	if err := mageutil.HasExec(tsc); err != nil {
		return err
	}
	if newer {
		return sh.Run(tsc, "-p", filepath.Join(webContentDir, "tsconfig.json"))
	}
	return nil
}

func Zip() error {
	mg.Deps(Content)
	return mageutil.ZipDir(webContentOutDir, webContentZip)
}
