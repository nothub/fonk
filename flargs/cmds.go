package flargs

import (
	"errors"
	"flag"
	"io"
)

// some common flags:
//   --help
//   --version
//   --silent
//   --verbose

var AlreadyParsedErr = errors.New("command already parsed")

type Cmd struct {
	Name string
	Desc string

	Flags []FlagMeta
	Cmds  []Cmd
	Args  []string

	SortFlags bool

	Usage func() string

	parsed        bool
	minArgs       uint8
	maxArgs       uint8
	output        io.Writer // nil means stderr; use Output() accessor
	errorHandling flag.ErrorHandling
}

func (cmd Cmd) Parse(args []string) error {
	if cmd.parsed {
		return AlreadyParsedErr
	}
	cmd.parsed = true

	cmd.Args = make([]string, len(args))
	copy(cmd.Args, args)

	// TODO: parse args to flags and subcommands recursively

	for _, child := range cmd.Cmds {
		err := child.Parse(args)
		if err != nil {
			return err
		}
	}

	return cmd.Validate()
}

func (cmd Cmd) Validate() error {
	// TODO: validate flags
	for _, child := range cmd.Cmds {
		err := child.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd Cmd) SetOutput(output io.Writer) {
	cmd.output = output
	for _, child := range cmd.Cmds {
		child.SetOutput(output)
	}
}
