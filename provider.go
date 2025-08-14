package time

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type provider struct{}

func NewProvider() flam.Provider {
	return &provider{}
}

func (provider) Id() string {
	return providerId
}

func (provider) Register(
	container *dig.Container,
) error {
	var e error
	provide := func(constructor any, opts ...dig.ProvideOption) bool {
		e := container.Provide(constructor, opts...)
		return e == nil
	}

	_ = provide(newFactory) &&
		provide(newFacade)

	return e
}
