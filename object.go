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

// Flags returns the flags with the given name in order.
func (objs Objects) Flags(name string) []ObjectFlag {
	var objFlags []ObjectFlag

	for _, obj := range objs {
		// Ensure the current object is a flag.
		objFlag, ok := obj.(ObjectFlag)
		if !ok {
			continue
		}

		// Check if the name is correct and append.
		if objFlag.Name == name {
			objFlags = append(objFlags, objFlag)
		}
	}

	return objFlags
}

// Flag returns the first flag with the given name.
//
// An error will be returned if there are no flags found or if there is more
// than one flag found.
func (objs Objects) Flag(name string) (ObjectFlag, error) {
	objFlags := objs.Flags(name)
	if len(objFlags) == 0 {
		return ObjectFlag{}, errFlagNotFound(name)
	}
	if len(objFlags) > 1 {
		return ObjectFlag{}, errTooManyFlags(name)
	}
	return objFlags[0], nil
}
