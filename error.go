package mojo

import "fmt"

// errInvalidFlag returns an invalid flag error.
func errInvalidFlag(name string) error {
	return fmt.Errorf("mojo: invalid flag %s", name)
}
