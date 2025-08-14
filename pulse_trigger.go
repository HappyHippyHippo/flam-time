package time

import (
	"time"
)

func newPulseTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	if callback == nil {
		return nil, newErrNilReference("callback")
	}

	timer := time.NewTimer(delay)
	closeCh := make(chan struct{})

	var t *trigger
	t = &trigger{
		delay:    delay,
		isClosed: false,
		closer: func() error {
			t.isClosed = true
			if timer != nil {
				closeCh <- struct{}{}
			}
			return nil
		},
		cleaner: func() error {
			if timer != nil {
				timer.Stop()
				timer = nil
				close(closeCh)
			}
			return nil
		},
	}

	go func(t *trigger) {
		if timer != nil {
			select {
			case <-timer.C:
				_ = callback()
			case <-closeCh:
			}
		}
		_ = t.cleaner()
	}(t)

	return t, nil
}
