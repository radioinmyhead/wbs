package main

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/fsnotify.v1"

	"path/filepath"
)

// WbsWatcher file wather struct
type WbsWatcher struct {
	w             *fsnotify.Watcher
	TargetDirs    []string
	ExcludeDirs   []string
	TargetFileExt []string
}

// initWatcher add watch target files to watcher
func (w *WbsWatcher) initWatcher() {
	var excludeDirRegexps []*regexp.Regexp
	for _, excludeDirStr := range w.ExcludeDirs {
		r := regexp.MustCompile(excludeDirStr)
		excludeDirRegexps = append(excludeDirRegexps, r)
	}

	for _, targetDir := range w.TargetDirs {
		filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
			for _, e := range excludeDirRegexps {
				if e.MatchString(path) {
					return nil
				}
			}
			for _, s := range w.TargetFileExt {
				if filepath.Ext(path) == s {
					log.Printf("start watching %s", path)
					err := w.w.Add(path)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			return nil
		})
	}
}

// NewFileWatcher create target file watcher
func NewWbsWatcher(config *WbsConfig) (*WbsWatcher, error) {
	var watcher *WbsWatcher
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("failed to create watcher: %s", err)
		return watcher, err
	}
	defer w.Close()
	watcher = &WbsWatcher{
		w:             w,
		TargetDirs:    config.WatchTargetDirs,
		ExcludeDirs:   config.WatchExcludeDirs,
		TargetFileExt: config.WatchFileExt,
	}
	watcher.initWatcher()
	return watcher, nil
}
