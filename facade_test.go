package time

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

func Test_Facade_TimeFunctions(t *testing.T) {
	t.Run("should parse valid duration string (ParseDuration)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.ParseDuration("1s")
			assert.Equal(t, time.Second, got)
			assert.NoError(t, e)
		}))
	})

	t.Run("should return current time (Now)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.WithinDuration(t, time.Now(), facade.Now(), 10*time.Millisecond)
		}))
	})

	t.Run("should return time since a point (Since)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			then := time.Now().Add(-5 * time.Second)
			assert.GreaterOrEqual(t, facade.Since(then), 5*time.Second)
		}))
	})

	t.Run("should return time until a point (Until)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			future := time.Now().Add(5 * time.Second)
			assert.GreaterOrEqual(t, facade.Until(future), 4*time.Second)
			assert.LessOrEqual(t, facade.Until(future), 5*time.Second)
		}))
	})

	t.Run("should return a location (FixedZone)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got := facade.FixedZone("test", 3600)
			assert.NotNil(t, got)
			assert.Equal(t, "test", got.String())
		}))
	})

	t.Run("should load a location (LoadLocation)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.LoadLocation("UTC")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			assert.Equal(t, "UTC", got.String())
		}))
	})

	t.Run("should create a date (Date)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			loc, _ := time.LoadLocation("UTC")
			got := facade.Date(2024, time.January, 1, 10, 30, 0, 0, loc)
			assert.Equal(t, 2024, got.Year())
			assert.Equal(t, time.January, got.Month())
		}))
	})

	t.Run("should parse a time string (Parse)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.Parse(time.RFC3339, "2024-01-01T10:30:00Z")
			assert.Equal(t, 2024, got.Year())
			assert.NoError(t, e)
		}))
	})

	t.Run("should parse a time string in a location (ParseInLocation)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			loc, _ := time.LoadLocation("UTC")
			got, e := facade.ParseInLocation(time.RFC3339, "2024-01-01T10:30:00Z", loc)
			assert.Equal(t, 2024, got.Year())
			assert.NoError(t, e)
		}))
	})

	t.Run("should create time from unix timestamp (Unix)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			ts := time.Now().Unix()
			assert.Equal(t, ts, facade.Unix(ts, 0).Unix())
		}))
	})

	t.Run("should create time from unix micro timestamp (UnixMicro)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			ts := time.Now().UnixMicro()
			assert.Equal(t, ts, facade.UnixMicro(ts).UnixMicro())
		}))
	})

	t.Run("should create time from unix milli timestamp (UnixMilli)", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			ts := time.Now().UnixMilli()
			assert.Equal(t, ts, facade.UnixMilli(ts).UnixMilli())
		}))
	})
}

func Test_Facade_NewPulseTrigger(t *testing.T) {
	t.Run("should return an error when callback is nil", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(1*time.Second, nil)
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrNilReference)
		}))
	})

	t.Run("should create a new pulse trigger", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		handler := func() error {
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(1*time.Second, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.NoError(t, got.Close())
		}))
	})

	t.Run("should return the defined delay", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		expectedDelay := 20 * time.Millisecond
		handler := func() error {
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(expectedDelay, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, expectedDelay, got.Delay())

			assert.NoError(t, got.Close())
		}))
	})
	t.Run("should return the correct closed state", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		handler := func() error {
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(20*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.False(t, got.IsClosed())

			assert.NoError(t, got.Close())
			assert.True(t, got.IsClosed())
		}))
	})

	t.Run("should execute the callback after the delay", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		var wg sync.WaitGroup
		wg.Add(1)

		callbackExecuted := false
		handler := func() error {
			callbackExecuted = true
			wg.Done()
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(10*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			wg.Wait()
			assert.True(t, callbackExecuted)
			assert.NoError(t, got.Close())
		}))
	})

	t.Run("should not execute the callback if closed before delay", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		callbackExecuted := false
		handler := func() error {
			callbackExecuted = true
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewPulseTrigger(20*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.NoError(t, got.Close())

			time.Sleep(30 * time.Millisecond)
			assert.False(t, callbackExecuted)
		}))
	})
}

func Test_Facade_NewRecurringTrigger(t *testing.T) {
	t.Run("should return an error when callback is nil", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(1*time.Second, nil)
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrNilReference)
		}))
	})

	t.Run("should create a new recurring trigger", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		handler := func() error {
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(1*time.Second, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.NoError(t, got.Close())
		}))
	})

	t.Run("should return the defined delay", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		handler := func() error {
			return nil
		}

		expectedDelay := 20 * time.Millisecond
		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(expectedDelay, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.Equal(t, expectedDelay, got.Delay())

			assert.NoError(t, got.Close())
		}))
	})
	t.Run("should return the correct closed state", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		handler := func() error {
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(20*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.False(t, got.IsClosed())

			assert.NoError(t, got.Close())
			assert.True(t, got.IsClosed())
		}))
	})

	t.Run("should execute the callback multiple times", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		var wg sync.WaitGroup
		wg.Add(3)

		callCount := 0
		handler := func() error {
			callCount++
			wg.Done()
			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(10*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			wg.Wait()
			assert.GreaterOrEqual(t, callCount, 3)

			assert.NoError(t, got.Close())
		}))
	})

	t.Run("should stop when callback returns an error", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		callCount := 0
		handler := func() error {
			callCount++
			if callCount == 2 {
				return errors.New("stop error")
			}

			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(10*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			time.Sleep(50 * time.Millisecond)
			assert.Equal(t, 2, callCount)

			assert.NoError(t, got.Close())
		}))
	})

	t.Run("should not execute the callback if closed", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		callCount := 0
		handler := func() error {
			callCount++

			return nil
		}

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.NewRecurringTrigger(10*time.Millisecond, handler)
			require.NotNil(t, got)
			require.NoError(t, e)

			time.Sleep(15 * time.Millisecond) // Allow it to run at least once
			assert.NoError(t, got.Close())

			firstCount := callCount
			time.Sleep(30 * time.Millisecond)
			assert.Equal(t, firstCount, callCount)
		}))
	})
}
