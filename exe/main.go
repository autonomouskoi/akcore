package exe

import (
	"context"
	"errors"
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
	"github.com/autonomouskoi/akcore/internal"
	"github.com/autonomouskoi/akcore/modules"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	"github.com/autonomouskoi/akcore/svc/log"
	svc "github.com/autonomouskoi/akcore/svc/pb"
	"github.com/autonomouskoi/akcore/web"
)

// Set up the internal dependencies for the bot. The app/ package has the true
// main() and will pull in modules for side-effects. This allows building custom
// version of AK just by creating a different main file with desired imports.

func Main() {
	// Trigger safe shutdown on common signals.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// set up the system tray GUI features
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

		// watch for menu events
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

		// launch the program proper and wait for it to return
		mainIsh(ctx, setStatus)
		systray.Quit()
	}

	// launch the system tray applet
	systray.Run(onReady, cancel)
}

func mainIsh(ctx context.Context, setStatus func(string)) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// figure out where our app path is. This is platform-dependent.
	appPath, err := run.AppPath()
	if err != nil {
		setStatus("Error determining app path: " + err.Error())
		<-ctx.Done()
		return
	}
	akCorePath := filepath.Join(appPath, "akcore")

	// provide the option to open AK's data folder
	mOpenDir := systray.AddMenuItem("Open data folder", "Open the folder with AK data")
	mOpenDir.Enable()
	go func() {
		for range mOpenDir.ClickedCh {
			run.ShowFolder(appPath)
		}
	}()

	// set up logging. Use a log file named for the date.
	// This is an initial logger at INFO level until we have a config
	logDir := filepath.Join(akCorePath, "logs")
	if err := os.MkdirAll(logDir, 0700); err != nil {
		setStatus("creating logs dir: " + err.Error())
		<-ctx.Done()
		return
	}
	initLogger, err := log.New(logDir, &svc.Config{})
	if err != nil {
		setStatus("Error creating initial logger: " + err.Error())
		<-ctx.Done()
		return
	}
	defer initLogger.Close()
	mainLog := initLogger.NewForSource("main")

	mainLog.Info("staring", "version", "v"+akcore.Version)

	// initialize the bus
	bus := bus.New(ctx)

	// initialize key-value storage
	kvPath := filepath.Join(akCorePath, "kv")
	kv, err := kv.New(kvPath)
	if err != nil {
		mainLog.Error("opening kv storage", "kvPath", kvPath, "error", err.Error())
		return
	}
	mainLog.Info("created kv storage", "kvPath", kvPath)

	// retrieve our config
	cfg, err := internal.GetConfig(*kv.WithPrefix([8]byte{}))
	if err != nil {
		mainLog.Error("getting config", "error", err.Error())
		return
	}
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = "localhost:8011"
	}

	// create our permanent logger
	initLogger.Close()
	masterLogger, err := log.New(logDir, cfg)
	if err != nil {
		setStatus("Error creating master logger: " + err.Error())
		<-ctx.Done()
		return
	}
	defer masterLogger.Close()
	mainLog = initLogger.NewForSource("main")

	// intialize the cache dir
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		mainLog.Error("getting user cache path", "error", err.Error())
		return
	}
	cacheDir = filepath.Join(cacheDir, "AutonomousKoi")

	eg, ctx := errgroup.WithContext(ctx)

	deps := &modutil.Deps{
		Bus:         bus,
		Log:         masterLogger,
		KV:          kv,
		CachePath:   cacheDir,
		StoragePath: appPath,
		Config:      cfg,
	}

	// launch the internal module, managing config, etc
	eg.Go(func() error { return internal.Start(ctx, deps) })

	// initialize the web service
	web := web.New("/", deps)
	deps.Web = web

	// start enabled modules
	eg.Go(func() error {
		return modules.Start(ctx, deps)
	})

	// initialize and start our web service
	server := &http.Server{
		Addr:    cfg.ListenAddress,
		Handler: web,
	}
	mainLog.Info("starting HTTP listener", "addr", cfg.ListenAddress)
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

	// indicate that we're running, wait for our errgroup to return
	setStatus("Running")
	if err := eg.Wait(); err != nil {
		mainLog.Error("in errgroup", "error", err.Error())
	}
	// safely close the KV store
	if err := kv.Close(); err != nil {
		mainLog.Error("closing kv storage", "error", err.Error())
	}
	mainLog.Info("shutting down")
}
