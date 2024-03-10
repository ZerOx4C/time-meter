package main

import (
	"path/filepath"
	"time-meter/util"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher       *fsnotify.Watcher
	filename      string
	busy          bool
	onFileChanged util.EventHandler
}

func (fw *FileWatcher) Initialize() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	fw.watcher = watcher

	return nil
}

func (fw *FileWatcher) Watch() {
	if fw.busy {
		panic("invalid operation.")
	}

	fullpath, err := filepath.Abs(fw.filename)
	if err != nil {
		panic(err)
	}

	fw.watcher.Add(filepath.Dir(fullpath))

	go func() {
		for {
			select {
			case event := <-fw.watcher.Events:
				if event.Name == fullpath {
					fw.onFileChanged.Invoke()
				}

			case <-fw.watcher.Errors:
			}
		}
	}()
}

func (fw *FileWatcher) Finalize() error {
	if err := fw.watcher.Close(); err != nil {
		return err
	}

	fw.watcher = nil

	return nil
}
