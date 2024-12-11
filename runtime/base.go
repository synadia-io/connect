package runtime

import (
	"context"
	"sync"
)

type WorkloadFoundation struct {
	rmut sync.Mutex

	running bool
	cancel  context.CancelFunc
}

func (w *WorkloadFoundation) Started(ctx context.Context) context.Context {
	w.rmut.Lock()
	defer w.rmut.Unlock()

	result, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.running = true
	return result
}

func (w *WorkloadFoundation) Stopped() {
	w.rmut.Lock()
	defer w.rmut.Unlock()
	w.cancel = nil
	w.running = false
}

func (w *WorkloadFoundation) IsRunning() bool {
	w.rmut.Lock()
	defer w.rmut.Unlock()
	return w.running
}

func (w *WorkloadFoundation) Stop() {
	w.rmut.Lock()
	defer w.rmut.Unlock()
	if w.cancel != nil {
		w.cancel()
	}
}
