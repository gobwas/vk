package syncutil

import (
	"errors"
	"time"
)

var ErrClosed = errors.New("limiter closed")

type Limiter struct {
	done    chan struct{}
	tickets chan struct{}
}

func NewLimiter(interval time.Duration, count int) *Limiter {
	l := &Limiter{
		done:    make(chan struct{}),
		tickets: make(chan struct{}, count),
	}
	fill := func() {
		for i := 0; i < count; i++ {
			select {
			case l.tickets <- struct{}{}:
			default:
				return
			}
		}
	}

	fill()

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-l.done:
				ticker.Stop()
				return
			case <-ticker.C:
				fill()
			}
		}
	}()

	return l
}

func (l *Limiter) Do(fn func()) error {
	select {
	case <-l.tickets:
	case <-l.done:
		return ErrClosed
	}
	fn()
	return nil
}

func (l *Limiter) Close() {
	close(l.done)
}
