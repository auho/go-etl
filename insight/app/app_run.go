package app

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

type Run struct {
	names []string
	fns   []func() error

	preFun   []func() error
	commands []*cobra.Command
}

// AddPreFunE
// 需要在 Run 或 RunE 之前执行的 func
func (r *Run) AddPreFunE(fn ...func() error) {
	r.preFun = append(r.preFun, fn...)
}

// RunPreFunE
// 在执行 Run 或 RunE 之前执行
func (r *Run) RunPreFunE(parentCmd *cobra.Command) error {
	fmt.Println(fmt.Sprintf("parent cmd[%s] run pre", parentCmd.Use))

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

func (r *Run) AddCommand(cs ...*cobra.Command) {
	for _, cmd := range cs {
		r.commands = append(r.commands, cmd)
	}
}

func (r *Run) RunCommandE(parentCmd *cobra.Command, args []string) error {
	var err error

	if parentCmd == nil {
		parentCmd = &cobra.Command{Use: "unknown"}
	}

	// 循环嵌套问题
	// 先 shallow clone，然后在清空，后续不会出现循环嵌套
	_commands := slices.Clone(r.commands)
	r.commands = nil

	fmt.Println(fmt.Sprintf("parent cmd[%s] start", parentCmd.Use))
	for _, _cmd := range _commands {
		fmt.Println(fmt.Sprintf("cmd[%s] begin ...", _cmd.Use))

		err = _cmd.RunE(_cmd, args)
		if err != nil {
			return fmt.Errorf("%s run error; %w", _cmd.Use, err)
		}

		fmt.Println(fmt.Sprintf("cmd[%s] end", _cmd.Use))
	}

	fmt.Println()
	fmt.Println(fmt.Sprintf("parent cmd[%s]:", parentCmd.Use))
	for _, _cmd := range _commands {
		fmt.Println(_cmd.Use)
	}

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
