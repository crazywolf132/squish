package watcher

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"squish/internal/utils"

	"squish/pkg/esbuild"
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

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					utils.Log("Modified file:", event.Name)
					if err := w.bundler.Bundle(); err != nil {
						utils.Log("Error bundling:", err)
					}
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

	utils.Log("Watching for changes in", w.srcDir)
	<-done
	return nil
}
