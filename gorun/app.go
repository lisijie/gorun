package gorun

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	defaultWatchExtensions = ".go,.toml,.ini,.yml"
)

type Config struct {
	Name             string            `toml:"name"`
	AppPath          string            `toml:"app_path"`
	WatchExcludeDirs string            `toml:"watch_exclude_dirs"`
	WatchExtensions  string            `toml:"watch_extensions"`
	BuildCommand     string            `toml:"build_cmd"`
	RunCommand       string            `toml:"run_cmd"`
	Environ          map[string]string `toml:"env"`
}

type App struct {
	cfg     *Config
	watcher *fsnotify.Watcher
	process *os.Process
	log     Logger
}

func New(cfg *Config) *App {
	if cfg.AppPath == "" {
		cfg.AppPath, _ = os.Getwd()
	} else {
		cfg.AppPath, _ = filepath.Abs(cfg.AppPath)
	}
	if cfg.WatchExtensions == "" {
		cfg.WatchExtensions = defaultWatchExtensions
	}
	app := &App{
		cfg: cfg,
		log: NullLogger,
	}
	return app
}

func (app *App) SetLogger(l Logger) {
	app.log = l
}

func (app *App) Run() (err error) {
	app.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}

	if err = app.initWatchDirs(); err != nil {
		return
	}

	app.buildAndRun()
	defer app.kill()

	exts := make(map[string]struct{})
	for _, v := range strings.Split(app.cfg.WatchExtensions, ",") {
		exts[strings.TrimSpace(v)] = struct{}{}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case e := <-app.watcher.Events:
			if _, ok := exts[filepath.Ext(e.Name)]; ok {
				app.buildAndRun()
			}
		case <-sig:
			return
		}
	}
}

func (app *App) initWatchDirs() error {
	excludeDirs := make(map[string]bool)
	if app.cfg.WatchExcludeDirs != "" {
		for _, v := range strings.Split(app.cfg.WatchExcludeDirs, ",") {
			absPath, err := filepath.Abs(filepath.Join(app.cfg.AppPath, v))
			if err == nil {
				excludeDirs[absPath] = true
			}
		}
	}
	err := filepath.Walk(app.cfg.AppPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name()[0] == '.' || excludeDirs[path] {
				return filepath.SkipDir
			}
			app.log.Debug("watch dir: ", path)
			return app.watcher.Add(path)
		}
		return nil
	})
	return err
}

func (app *App) buildAndRun() {
	app.log.Debug("----buildAndRun----")
	app.kill()
	if err := app.build(); err != nil {
		app.log.Error("build error: ", err)
		return
	}
	if err := app.run(); err != nil {
		app.log.Error("run error: ", err)
	}
}

func (app *App) build() error {
	var (
		errBuf bytes.Buffer
	)
	app.log.Debug("build app")
	cmd := exec.Command("sh", "-c", app.cfg.BuildCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil && errBuf.Len() > 0 {
		err = errors.New(errBuf.String())
	}
	return err
}

func (app *App) run() error {
	app.log.Debug("run app")
	cmd := exec.Command(app.cfg.RunCommand)
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	env := os.Environ()
	if len(app.cfg.Environ) > 0 {
		for k, v := range app.cfg.Environ {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	cmd.Env = env
	err := cmd.Start()
	if err != nil {
		return err
	}
	app.process = cmd.Process
	go func() {
		if err := cmd.Wait(); err != nil {
			app.log.Info("app exit: ", err)
		}
		app.process = nil
	}()
	return nil
}

func (app *App) kill() {
	if app.process != nil {
		app.log.Debug("kill app")
		app.process.Kill()
	}
}
