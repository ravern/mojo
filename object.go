package mojo

// Objects is a list of objects which represents some parsed arguments.
type Objects []Object

// Object represents a command, flag or argument.
type Object interface {
	object()
}

// ObjectCommand represents a command that has been parsed.
type ObjectCommand struct {
	Name string
}

func (ObjectCommand) object() {}

// ObjectFlag represents a flag that has been parsed.
type ObjectFlag struct {
	Name  string
	Value string

	// Bool indicates whether this flag was a bool flag.
	//
	// This means that the flag was passed without a value.
	Bool bool

	// MultipleFlagsStart indicates whether this flag was the start of
	// multiple flags (e.g. ls -al).
	//
	// Check AllowMultipleFlags in Config for more information.
	MultipleFlagsStart bool

	// MultipleFlagsEnd indicates whether this flag was the end of multiple
	// flags (e.g. ls -al).
	//
	// Check AllowMultipleFlags in Config for more information.
	MultipleFlagsEnd bool

	// CombinedFlagValues indicates whether this flag was combined
	// (i.e. --flag=value).
	//
	// Check DisallowCombinedFlagValues in Config for more information.
	CombinedFlagValues bool
}

func (ObjectFlag) object() {}

// ObjectArgument represents an argument that has been parsed.
type ObjectArgument struct {
	Value string
}

func (ObjectArgument) object() {}
