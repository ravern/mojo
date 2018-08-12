package mojo

// Config contains configuration that defines how to parse certain objects.
type Config struct {
	Root ConfigCommand

	// AllowUnconfiguredFlags indicates whether unconfigured flags are
	// allowed.
	//
	// If it isn't allowed, then unconfigured flags will result in an
	// invalid flag error.
	AllowUnconfiguredFlags bool

	// AllowMultipleFlags indicates whether combining multiple flags
	// (e.g. ls -al) is allowed.
	//
	// If it is allowed, such flags will be parsed as multiple flags with
	// the last one containing a value unless indicated otherwise by
	// configuration. If it isn't allowed, they will be parsed as a single
	// flag instead.
	AllowMutipleFlags bool

	// DisallowCombinedFlagValues indicates whether combining the flag and
	// value together (e.g. --flag=value) is not allowed.
	//
	// If it isn't allowed, then these flags will be taken literally, with
	// their name set to whatever the combined flag and value is.
	DisallowCombinedFlagValues bool

	// DisallowDoubleDash indicates whether the double dash (i.e. --) is
	// not allowed.
	//
	// If it is allowed, it will be parsed as an argument. However, if it
	// isn't allowed, then an error will occur since it will be treated
	// as a flag without a name.
	DisallowDoubleDash bool
}

// ConfigCommand contains configuration for a command.
type ConfigCommand struct {
	Name     string
	Commands []ConfigCommand
	Flags    []ConfigFlag
}

// ConfigFlag contains configuration for a flag.
type ConfigFlag struct {
	Name string
	Bool bool
}

// Command returns the command configuration for the command of the given name.
func (c ConfigCommand) Command(name string) (ConfigCommand, bool) {
	for _, cmd := range c.Commands {
		if cmd.Name == name {
			return cmd, true
		}
	}
	return ConfigCommand{}, false
}

// Flag returns the flag configuration for the command of the given name.
func (c ConfigCommand) Flag(name string) (ConfigFlag, bool) {
	for _, flag := range c.Flags {
		if flag.Name == name {
			return flag, true
		}
	}
	return ConfigFlag{}, false
}
