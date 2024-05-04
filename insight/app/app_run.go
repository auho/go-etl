package app

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

type Run struct {
	names []string
	fns   []func() error

	preFun []func() error
}

// AddPreRunE
// 需要在 Run 或 RunE 之前执行的 func
func (r *Run) AddPreRunE(fn ...func() error) {
	r.preFun = append(r.preFun, fn...)
}

// RunPreRunE
// 在执行 Run 或 RunE 之前执行
func (r *Run) RunPreRunE(cmd *cobra.Command) error {
	fmt.Println(fmt.Sprintf("cmd[%s] run pre", cmd.Use))

	_fns := slices.Clone(r.preFun)
	r.preFun = nil

	for _, fn := range _fns {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Run) RunCommandE(parentCmd *cobra.Command, args []string, cs ...*cobra.Command) error {
	return r.execCommandsE(parentCmd, args, cs)
}

func (r *Run) RunCommandsE(parentCmd *cobra.Command, args []string, cs ...[]*cobra.Command) error {
	var _cs []*cobra.Command
	for _, c := range cs {
		_cs = append(_cs, c...)
	}

	return r.execCommandsE(parentCmd, args, _cs)
}

func (r *Run) execCommandsE(parentCmd *cobra.Command, args []string, cs ...[]*cobra.Command) error {
	var err error

	var _commands []*cobra.Command
	for _, c := range cs {
		_commands = append(_commands, c...)
	}

	if parentCmd == nil {
		parentCmd = &cobra.Command{Use: "unknown"}
	}

	fmt.Println(fmt.Sprintf("parent cmd[%s] start", parentCmd.Use))
	for _, _cmd := range _commands {
		fmt.Println(fmt.Sprintf("cmd[%s] begin ...", _cmd.Use))

		err = _cmd.RunE(_cmd, args)
		if err != nil {
			return fmt.Errorf("%s run error; %w", _cmd.Use, err)
		}

		fmt.Println(fmt.Sprintf("cmd[%s] end", _cmd.Use))
	}

	fmt.Println("======")
	fmt.Println(fmt.Sprintf("parent cmd[%s]:", parentCmd.Use))
	for _, _cmd := range _commands {
		fmt.Println("  " + _cmd.Use)
	}
	fmt.Println()

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
