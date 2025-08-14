package time

import (
	"time"
)

func newRecurringTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	if callback == nil {
		return nil, newErrNilReference("callback")
	}

	timer := time.NewTicker(delay)
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
		for {
			select {
			case <-timer.C:
				if e := callback(); e != nil {
					_ = t.cleaner()
					return
				}
			case <-closeCh:
				_ = t.cleaner()
				return
			}
		}
	}(t)

	return t, nil
}
