package flock

import (
	"context"
	"sync"
	"time"

	extflock "github.com/gofrs/flock"
)

type FileLock struct {
	flk   *extflock.Flock
	mu    sync.Mutex
	owned bool
}

func New(path string) *FileLock {
	return &FileLock{
		flk: extflock.New(path),
	}
}

func (l *FileLock) Lock(ctx context.Context) (bool, error) {
	if l.owned {
		return true, nil
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	retryCh := make(chan bool, 1)
	retryCh <- true
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-ticker.C:
			retryCh <- true
		case <-retryCh:
			locked, err := l.flk.TryLock()
			if err != nil {
				return false, err
			}
			if locked {
				l.owned = true
				return true, nil
			}
		}
	}
}

func (l *FileLock) Unlock() error {
	if !l.owned {
		return nil
	}

	return l.flk.Unlock()
}
