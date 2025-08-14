package time

import (
	flam "github.com/happyhippyhippo/flam"
)

func newErrNilReference(
	arg string,
) error {
	return flam.NewErrorFrom(
		flam.ErrNilReference,
		arg)
}
