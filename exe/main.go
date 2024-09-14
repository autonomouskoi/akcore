package exe

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"fyne.io/systray"
	"github.com/pkg/browser"
	"golang.org/x/sync/errgroup"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/exe/run"
	"github.com/autonomouskoi/akcore/modules"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	"github.com/autonomouskoi/akcore/web"
)

func Main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	onReady := func() {
		// systray.SetTitle("AutonomousKoi")
		systray.SetIcon(run.IconBytes)
		systray.SetTooltip("The AutomousKoi Bot v" + akcore.Version)

		mBrowse := systray.AddMenuItem("Controls", "Open AK's controls in your browser")
		mBrowse.Enable()

		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		mQuit.Enable()

		systray.AddSeparator()

		mStatus := systray.AddMenuItem("Status: Starting...", "Current program status")
		setStatus := func(status string) {
			mStatus.SetTitle("Status: " + status)
		}

		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					setStatus("Shutting down...")
					cancel()
					return
				case <-mBrowse.ClickedCh:
					// TODO: handle error
					_ = browser.OpenURL("http://localhost:8011/")
				}
			}
		}()

		mainIsh(ctx, setStatus)
		systray.Quit()
	}

	systray.Run(onReady, cancel)
}

func mainIsh(ctx context.Context, setStatus func(string)) {
	appPath, err := run.AppPath()
	if err != nil {
		setStatus("Error determining app path: " + err.Error())
		<-ctx.Done()
		return
	}
	akCorePath := filepath.Join(appPath, "akcore")

	mOpenDir := systray.AddMenuItem("Open data folder", "Open the folder with AK data")
	mOpenDir.Enable()
	go func() {
		for range mOpenDir.ClickedCh {
			run.ShowFolder(akCorePath)
		}
	}()

	logDir := filepath.Join(akCorePath, "logs")
	if err := os.MkdirAll(logDir, 0700); err != nil {
		setStatus("Error creating logs folder: " + err.Error())
		<-ctx.Done()
		return
	}
	logFilePath := filepath.Join(logDir, time.Now().Format("ak-20060102.log"))

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		log.Fatal("error creating log file: ", err)
	}
	defer logFile.Close()

	bus := bus.New(ctx)
	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	log.Info("staring", "version", "v"+akcore.Version)

	kvPath := filepath.Join(akCorePath, "kv")
	kv, err := kv.New(kvPath)
	if err != nil {
		log.Error("opening kv storage", "kvPath", kvPath, "error", err.Error())
		os.Exit(-1)
	}
	log.Debug("created kv storage", "kvPath", kvPath)
	defer func() {
		if err := kv.Close(); err != nil {
			log.Error("closing kv storage", "error", err.Error())
		}
	}()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Error("getting user cache path", "error", err.Error())
		os.Exit(-1)
	}

	deps := &modutil.Deps{
		Bus:         bus,
		Log:         log,
		KV:          kv,
		CachePath:   cacheDir,
		StoragePath: akCorePath,
	}

	web := web.New("/", deps)
	deps.Web = web

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return modules.Start(ctx, deps)
	})

	server := &http.Server{
		Addr:    "localhost:8011",
		Handler: web,
	}
	log.Info("starting HTTP listener", "addr", "localhost:8011")
	eg.Go(func() error {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		<-ctx.Done()
		serverCtx, serverCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer serverCancel()
		return server.Shutdown(serverCtx)
	})

	setStatus("Running")
	if err := eg.Wait(); err != nil {
		log.Error("in errgroup", "error", err.Error())
	}
}
