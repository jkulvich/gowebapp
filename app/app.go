package app

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Config - WebApp config
type Config struct {
	Engine  string `yaml:"engine"`
	Address string `yaml:"addr"`
}

// App - WebApp instance
type App struct {
	ctx     context.Context
	lg      *logrus.Logger
	conf    *Config
	stopErr chan error
}

// NewApp - create new WebApp instance
func NewApp(ctx context.Context, lg *logrus.Logger, conf *Config) *App {
	return &App{
		ctx:  ctx,
		lg:   lg,
		conf: conf,
		stopErr: make(chan error, 1),
	}
}

// Run - run the WebApp process
func (app *App) Run() error {
	// web engine, one of ["chromium", "google-chrome"] or chrome like
	enginePath, err := exec.LookPath(app.conf.Engine)
	if err != nil {
		return err
	}

	// user data temp dir
	tmpDir, err := ioutil.TempDir(os.TempDir(), "gowebapp")
	if err != nil {
		return err
	}
	userDataDir := fmt.Sprintf("--user-data-dir=%s", tmpDir)

	// launch the site as app
	appAddr := fmt.Sprintf("--app=http://%s", app.conf.Address)

	// process arguments
	args := []string{
		enginePath,
		userDataDir,
		appAddr,
	}

	// trying to spawn the process
	app.lg.Infof("app starting as: %s", strings.Join(args, " "))
	proc, err := os.StartProcess(enginePath, args, &os.ProcAttr{})
	if err != nil {
		return err
	}
	app.lg.Infof("app started")

	// waiting for process error or closed
	go func() {
		for {
			st, err := proc.Wait()
			if err != nil {
				app.stopErr <- err
				break
			}
			if st.Exited() {
				app.stopErr <- nil
				break
			}
		}
	}()

	// waiting for process stop signal
	if err := <- app.stopErr; err != nil {
		return err
	}

	// stopping the web app process
	app.lg.Infof("app stopping")
	_ = proc.Kill() //< try to kill anyway
	app.lg.Infof("app stopped")

	// trying to remove temp dir
	if err := os.RemoveAll(tmpDir); err != nil {
		return err
	}

	return nil
}

// Close - stop the WebApp process
func (app *App) Close() {
	if len(app.stopErr) == 0 {
		app.stopErr <- nil
	}
}
