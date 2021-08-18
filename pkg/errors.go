package pkg

import "github.com/pkg/errors"

// Wrap wraps two given errors into one generic error.
func Wrap(err1, err2 error) error {
	return errors.Wrap(err1, err2.Error())
}
