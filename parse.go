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

		// TODO: Parse as flag
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

// configCommands returns the command configurations of the given command stack,
// with the root command being last.
//
// It is assumed that the command stack is safe. It it isn't, weird behaviour
// will occur due to usage of invalid return values (i.e. the ok flag is
// ignored).
func configCommands(conf Config, commands []string) []ConfigCommand {
	cmds := []ConfigCommand{conf.Root}
	for _, command := range commands {
		cmd, _ := cmds[0].Command(command)
		cmds = append([]ConfigCommand{cmd}, cmds...)
	}
	return cmds
}
