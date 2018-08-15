package mojo

import "fmt"

// Possible wrapped errors.
var (
	ErrInvalidFlag            = fmt.Errorf("mojo: invalid flag")
	ErrIncompleteMultipleFlag = fmt.Errorf("mojo: incomplete multiple flag")
	ErrFlagNotFound           = fmt.Errorf("mojo: flag not found")
	ErrArgumentNotFound       = fmt.Errorf("mojo: argument not found")

	// ErrUnconfiguredFlag occurs during parsing when a flag that does not
	// exist in the configuration is found.
	ErrUnconfiguredFlag = fmt.Errorf("mojo: unconfigured flag")

	// ErrUnexpectedArrayFlag occurs when more than one flag with the same name
	// is found when only one is requested.
	ErrUnexpectedArrayFlag = fmt.Errorf("mojo: unexpected array flag")
)

// FlagError represents a flag error.
type FlagError struct {
	Name string
	Err  error
}

func (err FlagError) Error() string {
	if err.Name == "" {
		return err.Err.Error()
	}
	return fmt.Sprintf("%v: %s", err.Err, err.Name)
}

// ArgumentError represents an argument error.
type ArgumentError struct {
	Index int
	Err   error
}

func (err ArgumentError) Error() string {
	return fmt.Sprintf("%v: %d", err.Err, err.Index)
}
