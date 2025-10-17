package time

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type provider struct{}

func NewProvider() flam.Provider {
	return &provider{}
}

func (*provider) Id() string {
	return providerId
}

func (*provider) Register(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	registerer := flam.NewRegisterer()
	registerer.Queue(newFactory)
	registerer.Queue(newFacade)

	return registerer.Run(container)
}
