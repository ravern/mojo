package mojo

import (
	"strings"
)

// Parse parses the given arguments into objects using the given configuration.
//
// The first argument given should be the name of the root command (e.g. git).
func Parse(conf Config, args []string) (Objects, error) {
	if len(args) < 1 {
		panic("runtime error: index out of bounds")
	}
	return parseCommand(conf, []string{}, args)
}

// parseCommand parses the given arguments into objects using the given
// configuration, in the context of the current command stack.
//
// The first argument given should be the name of the root command (e.g. git).
func parseCommand(conf Config, commands []string, args []string) (Objects, error) {
	var objs []Object

	// Append the root command to the objects and the command stack.
	objs = append(objs, ObjectCommand{Name: args[0]})
	commands = append(commands, args[0])
	args = args[1:]

	// Go through the rest of the arguments.
	for len(args) > 0 {
		// Determine if the argument is a command or argument.
		if !strings.HasPrefix(args[0], "-") {
			// Check for command.
			if _, ok := configCommands(conf, commands)[0].Command(args[0]); ok {
				// Parse the subcommand.
				subObjs, err := parseCommand(conf, commands, args)
				if err != nil {
					return nil, err
				}

				// Append everything and break, since parsing
				// is DONE!
				objs = append(objs, subObjs...)
				break
			}

			// Append as argument.
			objs = append(objs, ObjectArgument{Value: args[0]})
			args = args[1:]
			continue
		}

		// Check for the double dash only.
		if args[0] == "--" {
			obj, err := parseDoubleDash(conf)
			if err != nil {
				return nil, err
			}

			// Append the double dash.
			objs = append(objs, obj)
			args = args[1:]
			continue
		}

		// Check for combined flag value.
		//
		// The following code:
		// - Modifies `args`
		// - Does not modify `objs`
		var combinedFlagValue bool
		if i := strings.Index(args[0], "="); !conf.DisallowCombinedFlagValues && i != -1 {
			combinedFlagValue = true

			// Split into two different arguments and prepend them
			// back into the arguments.
			//
			// This means that ["--flag=value", "argument"] will
			// become ["--flag", "value", "argument"].
			args = append([]string{args[0][:i], args[0][i+1:]}, args[1:]...)
		}

		// Check for single dash flag with multiple characters.
		//
		// The following code:
		// - Modifies `args`
		// - Modifies `objs`
		var mutlipleFlagsEnd bool
		if conf.AllowMutipleFlags && !strings.HasPrefix(args[0], "--") && len(args[0]) > 3 {
			mutlipleFlagsEnd = true

			// Split the characters into individual flags.
			names := strings.Split(args[0][1:], "")
			for i := range names {
				names[i] = "-" + names[i]
			}

			// Remove the last flag in case it has a value.
			lastName := names[len(names)-1]
			names = names[:len(names)-1]

			// Add the individual flags as bools with the first
			// having the multiple flag start indication.
			for i, name := range names {
				obj, err := parseFlag(conf, commands, name, nil)
				if err != nil {
					return nil, err
				}

				if i == 0 {
					obj.MultipleFlagsStart = true
				}
				obj.CombinedFlagValues = combinedFlagValue

				objs = append(objs, obj)
				args = args[1:]
			}

			// Prepend the last flag back into arguments.
			//
			// This means that ["-abcd", "value"] will become
			// ["-d", "value"], with the bool flags "-a", "-b" and
			// "-c" already appended.
			args = append([]string{lastName}, args[1:]...)
		}

		// If there is a next value, and it isn't a flag, then set the
		// value.
		var value *string
		if len(args) > 1 && !strings.HasPrefix(args[1], "-") {
			tmp := string([]byte(args[1]))
			value = &tmp
		}

		// Parse the flag (FINALLY).
		obj, err := parseFlag(conf, commands, args[0], value)
		if err != nil {
			return nil, err
		}

		obj.CombinedFlagValues = combinedFlagValue
		obj.MultipleFlagsEnd = mutlipleFlagsEnd

		objs = append(objs, obj)
		args = args[1:]
		// If the flag isn't a bool flag, then remove the next argument
		// as well.
		if !obj.Bool {
			args = args[1:]
		}
	}

	return objs, nil
}

// parseDoubleDash parses a double dash argument based on the given
// configuration.
func parseDoubleDash(conf Config) (Object, error) {
	if conf.DisallowDoubleDash {
		return ObjectArgument{}, errInvalidFlag("--")
	}
	return ObjectArgument{
		Value: "--",
	}, nil
}

// parseFlag parses the given name and value as a flag based on the given
// configuration.
//
// Pass a value if the next argument is not a flag. Check if the value was used
// based on whether the resulting flag has Bool set.
func parseFlag(conf Config, commands []string, name string, value *string) (ObjectFlag, error) {
	// Find the flag in the configuration. If the configuration cannot be
	// found and unconfigured flags are not allowed, then return invalid
	// flag.
	flag, ok := configFlag(conf, commands, name)
	if !conf.AllowUnconfiguredFlags && !ok {
		return ObjectFlag{}, errInvalidFlag(name)
	}

	// If the flag is a not bool flag, but a value is not provided, then the
	// flag is invalid.
	if ok && !flag.Bool && value == nil {
		return ObjectFlag{}, errInvalidFlag(name)
	}

	// Create the flag and assign the value.
	obj := ObjectFlag{Name: name}
	if ok {
		if flag.Bool {
			obj.Bool = true
		} else {
			obj.Value = *value
		}
	} else {
		if value == nil {
			obj.Bool = true
		} else {
			obj.Value = *value
		}
	}

	return obj, nil
}

// configCommands returns the command configurations of the given command stack,
// with the root command being last.
//
// It is assumed that the command stack is safe. It it isn't, weird behaviour
// will occur due to usage of invalid return values (i.e. the ok flag is
// ignored).
func configCommands(conf Config, commands []string) []ConfigCommand {
	cmds := []ConfigCommand{conf.Root}
	for _, command := range commands[1:] {
		cmd, _ := cmds[0].Command(command)
		cmds = append([]ConfigCommand{cmd}, cmds...)
	}
	return cmds
}

// configFlag returns the flag configuration of the flag with the given name,
// with precedence given to configuration in the subcommands.
func configFlag(conf Config, commands []string, name string) (ConfigFlag, bool) {
	cmds := configCommands(conf, commands)
	for _, cmd := range cmds {
		if flag, ok := cmd.Flag(name); ok {
			return flag, ok
		}
	}
	return ConfigFlag{}, false
}
