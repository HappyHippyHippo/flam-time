package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.time.provider", NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return an error when a nil container is provided", func(t *testing.T) {
		assert.Error(t, NewProvider().Register(nil))
	})

	t.Run("should register the provider services", func(t *testing.T) {
		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NotNil(t, facade)
		}))
	})
}
