package time

import (
	"time"
)

type Facade interface {
	ParseDuration(s string) (time.Duration, error)
	Since(t time.Time) time.Duration
	Until(t time.Time) time.Duration
	FixedZone(name string, offset int) *time.Location
	LoadLocation(name string) (*time.Location, error)
	Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time
	Now() time.Time
	Parse(layout, value string) (time.Time, error)
	ParseInLocation(layout, value string, loc *time.Location) (time.Time, error)
	Unix(sec int64, nsec int64) time.Time
	UnixMicro(usec int64) time.Time
	UnixMilli(msec int64) time.Time

	NewPulseTrigger(delay time.Duration, callback Callback) (Trigger, error)
	NewRecurringTrigger(delay time.Duration, callback Callback) (Trigger, error)
}

type facade struct {
	factory *factory
}

func newFacade(
	factory *factory,
) Facade {
	return &facade{
		factory: factory,
	}
}

func (facade) ParseDuration(
	s string,
) (time.Duration, error) {
	return time.ParseDuration(s)
}

func (facade) Since(
	t time.Time,
) time.Duration {
	return time.Since(t)
}

func (facade) Until(
	t time.Time,
) time.Duration {
	return time.Until(t)
}

func (facade) FixedZone(
	name string,
	offset int,
) *time.Location {
	return time.FixedZone(name, offset)
}

func (facade) LoadLocation(
	name string,
) (*time.Location, error) {
	return time.LoadLocation(name)
}

func (facade) Date(
	year int,
	month time.Month,
	day,
	hour,
	min,
	sec,
	nsec int,
	loc *time.Location,
) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, loc)
}

func (facade) Now() time.Time {
	return time.Now()
}

func (facade) Parse(
	layout,
	value string,
) (time.Time, error) {
	return time.Parse(layout, value)
}

func (facade) ParseInLocation(
	layout,
	value string,
	loc *time.Location,
) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}

func (facade) Unix(
	sec int64,
	nsec int64,
) time.Time {
	return time.Unix(sec, nsec)
}

func (facade) UnixMicro(
	usec int64,
) time.Time {
	return time.UnixMicro(usec)
}

func (facade) UnixMilli(
	msec int64,
) time.Time {
	return time.UnixMilli(msec)
}

func (facade facade) NewPulseTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	return facade.factory.NewPulseTrigger(delay, callback)
}

func (facade facade) NewRecurringTrigger(
	delay time.Duration,
	callback Callback,
) (Trigger, error) {
	return facade.factory.NewRecurringTrigger(delay, callback)
}
