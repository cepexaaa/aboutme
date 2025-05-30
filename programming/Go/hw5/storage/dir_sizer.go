package storage

import (
	"context"
	"sync"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	wg              sync.WaitGroup
	mu              sync.Mutex
	size            int64
	count           int64
}

func NewSizer() DirSizer {
	return &sizer{maxWorkersCount: 20}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}
	errChan := make(chan error, 1)
	defer close(errChan)
	sem := make(chan struct{}, a.maxWorkersCount)
	defer close(sem)

	a.wg.Add(1)
	a.processDir(ctx, d, sem, errChan)
	a.wg.Wait()
	select {
	case err := <-errChan:
		return Result{}, err
	default:
		return Result{Size: a.size, Count: a.count}, nil
	}
}

func (a *sizer) processDir(ctx context.Context, d Dir, sem chan (struct{}), errChan chan (error)) {
	defer a.wg.Done()
	select {
	case <-ctx.Done():
		return
	default:
	}
	dirs, files, err := d.Ls(ctx)
	if err != nil {
		a.sendError(errChan, err)
		return
	}

	runInGoroutine(&a.wg, files, sem, func(file File) {
		defer a.wg.Done()
		size, err := file.Stat(ctx)
		if err != nil {
			a.sendError(errChan, err)
			return
		}
		a.updateResult(size)
	})

	runInGoroutine(&a.wg, dirs, sem, func(dir Dir) {
		a.processDir(ctx, dir, sem, errChan)
	})
}

func runInGoroutine[T any](wg *sync.WaitGroup, items []T, sem chan struct{}, processItem func(T)) {
	for _, item := range items {
		wg.Add(1)
		sem <- struct{}{}
		go func(item T) {
			defer func() { <-sem }()
			processItem(item)
		}(item)
	}
}

func (a *sizer) updateResult(size int64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.size += size
	a.count++
}

func (s *sizer) sendError(errCh chan error, err error) {
	select {
	case errCh <- err:
	default:
		// error already exists in errChan
	}
}
