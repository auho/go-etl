package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Run struct {
	names []string
	fns   []func() error

	commands []*cobra.Command
}

func (r *Run) AddCommand(cmd *cobra.Command) {
	r.commands = append(r.commands, cmd)
}

func (r *Run) RunCommandE() error {
	var err error
	for _, _cmd := range r.commands {
		err = _cmd.Execute()
		if err != nil {
			return err
		}
	}

	r.commands = nil

	return nil
}

func (r *Run) AddRunE(name string, fn func() error) *Run {
	r.names = append(r.names, name)
	r.fns = append(r.fns, fn)

	return r
}

func (r *Run) RunE() error {
	for i, fn := range r.fns {
		err := fn()
		if err != nil {
			return fmt.Errorf("%s; %w", r.names[i], err)
		}
	}

	r.fns = nil

	return nil
}
