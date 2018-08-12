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
}

func (ObjectFlag) object() {}

// ObjectArgument represents an argument that has been parsed.
type ObjectArgument struct {
	Value string
}

func (ObjectArgument) object() {}
