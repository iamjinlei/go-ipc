package flock_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/iamjinlei/go-ipc/flock"
)

func TestFlockBasic(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lk := flock.New("/tmp/test_flock_basic.lock")
	locked, err := lk.Lock(ctx)
	assert.NoError(t, err)
	assert.True(t, locked)
	assert.NoError(t, lk.Unlock())
}

func TestFlockExclusiveness(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lockPath := "/tmp/test_flock_excl.lock"

	cnt := atomic.NewInt64(0)
	lockedSignal := make(chan bool, 1)
	unlockSignal := make(chan bool)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			lk := flock.New(lockPath)
			locked, err := lk.Lock(ctx)
			assert.NoError(t, err)
			assert.True(t, locked)
			cnt.Inc()
			lockedSignal <- true
			<-unlockSignal
			require.NoError(t, lk.Unlock())
		}()
	}

	<-lockedSignal
	assert.Equal(t, int64(1), cnt.Load())
	unlockSignal <- true
	<-lockedSignal
	assert.Equal(t, int64(2), cnt.Load())
	unlockSignal <- true
	<-lockedSignal
	assert.Equal(t, int64(3), cnt.Load())
	unlockSignal <- true
	wg.Wait()
}
