package watcher

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"os/signal"
	"path/filepath"
	"squish/internal/utils"
	"squish/pkg/esbuild"
	"syscall"
	"time"
)

type Watcher struct {
	bundler *esbuild.Bundler
	srcDir  string
}

func NewWatcher(bundler *esbuild.Bundler, srcDir string) *Watcher {
	return &Watcher{
		bundler: bundler,
		srcDir:  srcDir,
	}
}

func (w *Watcher) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	var rebuildTimer *time.Timer
	debounceDuration := 100 * time.Millisecond

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if rebuildTimer != nil {
						rebuildTimer.Stop()
					}
					rebuildTimer = time.AfterFunc(debounceDuration, func() {
						utils.Log("Detected changes, rebuilding...")
						if err := w.bundler.Bundle(); err != nil {
							utils.Log("Error bundling:", err)
						}
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Log("Error:", err)
			}
		}
	}()

	err = filepath.Walk(w.srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Setup signal handling for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		utils.Log("Received interrupt signal, shutting down watcher...")
		close(done)
	}()

	utils.Log("Watching for changes in", w.srcDir)
	<-done
	return nil
}
