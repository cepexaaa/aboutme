package watcher

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type EventType string

const (
	EventTypeFileCreate  EventType = "file_created"
	EventTypeFileRemoved EventType = "file_removed"
)

var ErrDirNotExist = errors.New("dir does not exist")

type Event struct {
	Type EventType
	Path string
}

type Watcher struct {
	Events          chan Event
	refreshInterval time.Duration
	mu              sync.Mutex
	files           map[string]struct{}
}

func NewDirWatcher(refreshInterval time.Duration) *Watcher {
	return &Watcher{
		refreshInterval: refreshInterval,
		Events:          make(chan Event),
		files:           make(map[string]struct{}),
	}
}

func (w *Watcher) WatchDir(ctx context.Context, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrDirNotExist
	}

	err := w.initFiles(path)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(w.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.scanDirectory(path); err != nil {
				fmt.Printf("Scan error: %v", err)
				return err
			}
		}
	}
}

func (w *Watcher) initFiles(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.walkDirAndFiles(path, w.files)
}

func (w *Watcher) walkDirAndFiles(path string, dataFiles map[string]struct{}) error {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			dataFiles[filePath] = struct{}{}
		}
		return nil
	})
	return err
}

func (w *Watcher) updateDirInfo(currentFiles, lastFiles map[string]struct{}, message EventType) {
	for filePath := range currentFiles {
		if _, exists := lastFiles[filePath]; !exists {
			w.Events <- Event{Type: message, Path: filePath}
		}
	}
}

func (w *Watcher) findRemovedFiles(currentFiles map[string]struct{}) {
	w.updateDirInfo(w.files, currentFiles, EventTypeFileRemoved)
}

func (w *Watcher) findCreatedFiles(currentFiles map[string]struct{}) {
	w.updateDirInfo(currentFiles, w.files, EventTypeFileCreate)
}

func (w *Watcher) scanDirectory(path string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	currentFiles := make(map[string]struct{})
	err := w.walkDirAndFiles(path, currentFiles)
	if err != nil {
		return fmt.Errorf("error scanning directory: %w", err)
	}

	w.findCreatedFiles(currentFiles)
	w.findRemovedFiles(currentFiles)
	w.files = currentFiles
	return nil
}

func (w *Watcher) Close() {
	close(w.Events)
}
