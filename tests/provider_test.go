package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	time "github.com/happyhippyhippo/flam-time"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, time.NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.time.provider", time.NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	container := dig.New()
	require.NoError(t, time.NewProvider().Register(container))

	assert.NoError(t, container.Invoke(func(facade time.Facade) {
		assert.NotNil(t, facade)
	}))
}
