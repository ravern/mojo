package mojo

// Objects is a list of objects which represents some parsed arguments.
type Objects []Object

// Object represents a command, flag or argument.
type Object interface {
	object()
}

// CommandObject represents a command that has been parsed.
type CommandObject struct {
	Name string
}

func (CommandObject) object() {}

// FlagObject represents a flag that has been parsed.
type FlagObject struct {
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

func (FlagObject) object() {}

// ArgumentObject represents an argument that has been parsed.
type ArgumentObject struct {
	Value string
}

func (ArgumentObject) object() {}

// Flags returns the flags with the given name in order.
func (objs Objects) Flags(name string) []FlagObject {
	var flagObjs []FlagObject

	for _, obj := range objs {
		// Ensure the current object is a flag.
		flagObj, ok := obj.(FlagObject)
		if !ok {
			continue
		}

		// Check if the name is correct and append.
		if flagObj.Name == name {
			flagObjs = append(flagObjs, flagObj)
		}
	}

	return flagObjs
}

// Flag returns the first flag with the given name.
//
// An error will be returned if there are no flags found or if there is more
// than one flag found.
func (objs Objects) Flag(name string) (FlagObject, error) {
	flagObjs := objs.Flags(name)
	if len(flagObjs) == 0 {
		return FlagObject{}, FlagError{
			Name: name,
			Err:  ErrFlagNotFound,
		}
	}
	if len(flagObjs) > 1 {
		return FlagObject{}, FlagError{
			Name: name,
			Err:  ErrTooManyFlags,
		}
	}
	return flagObjs[0], nil
}

// Argument returns the argument at the given index.
func (objs Objects) Argument(i int) (ArgumentObject, error) {
	var j int

	for _, obj := range objs {
		argObj, ok := obj.(ArgumentObject)
		if !ok {
			continue
		}

		if j == i {
			return argObj, nil
		}

		j++
	}

	return ArgumentObject{}, ArgumentError{
		Index: i,
		Err:   ErrArgumentNotFound,
	}
}
