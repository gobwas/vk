package logutil

import (
	"container/ring"
	"io"
	"sync"
)

type RingLogger struct {
	mu   sync.Mutex
	ring *ring.Ring
}

func NewRingLogger(n int) *RingLogger {
	return &RingLogger{
		ring: ring.New(n),
	}
}

func (r *RingLogger) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	next := r.ring.Next()
	next.Value = string(p)
	r.ring = r.ring.Move(1)
	return len(p), nil
}

func (r *RingLogger) Interceptor() func(io.Writer) {
	return func(out io.Writer) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.ring.Do(func(v interface{}) {
			if v == nil {
				return
			}
			out.Write([]byte(v.(string)))
		})
	}
}
