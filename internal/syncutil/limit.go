package syncutil

import (
	"context"
	"time"
)

type Limiter struct {
	ctx     context.Context
	tickets chan struct{}
}

func NewLimiter(ctx context.Context, interval time.Duration, count int) *Limiter {
	l := &Limiter{
		ctx:     ctx,
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
			case <-ctx.Done():
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
	case <-l.ctx.Done():
		return l.ctx.Err()
	}

	fn()

	return nil
}
