package time

import (
	"time"
)

type factory struct{}

func newFactory() *factory {
	return &factory{}
}

func (factory *factory) NewPulseTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	return newPulseTrigger(delay, callback)
}

func (factory *factory) NewRecurringTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	return newRecurringTrigger(delay, callback)
}
