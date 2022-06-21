package helper

import "sync/atomic"

type Counter struct {
	finished uint64
}

func (c *Counter) Get() uint64 {
	return atomic.LoadUint64(&c.finished)
}

func (c *Counter) Inc() {
	atomic.AddUint64(&c.finished, 1)
}
