// Created by Hisen at 2022/3/1.
package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"runtime"
	"strings"
)
import (
	"os"
)

type Command struct {
	usage    string
	desc     string
	options  CliOptions
	commands []*Command
	runFunc  RunCommandFunc
}

type CommandsOption func(command *Command)

func WithCommandOptions(opt CliOptions) CommandsOption {
	return func(command *Command) {
		command.options = opt
	}
}

type RunCommandFunc func(args []string) error

func WithCommandRunFunc(commandFunc RunCommandFunc) CommandsOption {

	return func(command *Command) {
		command.runFunc = commandFunc
	}

}

func NewCommand(usage, desc string, opts ...CommandsOption) *Command {
	c := &Command{usage: usage, desc: desc}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Command) AddCommand(cmd *Command) {
	c.commands = append(c.commands, cmd)
}

func (c *Command) AddCommands(cmd ...*Command) {
	c.commands = append(c.commands, cmd...)
}

func (c *Command) cobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.usage,
		Short: c.desc,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stdout)

	if c.runFunc != nil {
		cmd.Run = c.runCommand
	}
	if c.options != nil {
		for _, f := range c.options.Flags().FlagSets {
			cmd.Flags().AddFlagSet(f)
		}
	}
	addHelpCommandFlag(c.usage, cmd.Flags())
	return cmd

}

func (c *Command) runCommand(cmd *cobra.Command, args []string) {
	if c.runFunc != nil {
		if err := c.runFunc(args); err != nil {
			fmt.Printf("%v %v\n", color.RedString("Error:"), err)
			os.Exit(1)
		}
	}
}
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}
