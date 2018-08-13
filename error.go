package mojo

import "fmt"

// errInvalidFlag returns an invalid flag error.
func errInvalidFlag(name string) error {
	return fmt.Errorf("mojo: invalid flag %s", name)
}

// errUnconfiguredFlag returns an unconfigured flag error.
//
// This error occurs when a flag that does not exist in the configuration is
// found in the arguments.
func errUnconfiguredFlag(name string) error {
	return fmt.Errorf("mojo: unconfigured flag %s", name)
}

// errIncompleteMultipleFlag returns an incomplete multiple flag error.
func errIncompleteMultipleFlag() error {
	return fmt.Errorf("mojo: incomplete multiple flag")
}
