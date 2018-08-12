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
			// TODO: Parse as command or argument.
			args = args[1:]
			continue
		}

		// Check for the double dash.
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
