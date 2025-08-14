package time

import (
	"io"
	"time"
)

type Trigger interface {
	io.Closer

	IsClosed() bool
	Delay() time.Duration
}

type trigger struct {
	delay    time.Duration
	isClosed bool
	closer   func() error
	cleaner  func() error
}

func (trigger *trigger) Close() error {
	return trigger.closer()
}

func (trigger *trigger) IsClosed() bool {
	return trigger.isClosed
}

func (trigger *trigger) Delay() time.Duration {
	return trigger.delay
}
