// Created by Hisen at 2022/3/1.
package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
)

const (
	flagHelp          = "help"
	flagHelpShortHand = "H"
)

func helpCommand(name string) *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "help about any command.",
		Long: `Help provides help for any command in the application.
Simply type ` + name + ` help [path to command] for full details.`,
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknow help topy %#q\n", args)
				_ = c.Root().Usage()
			} else {
				cmd.InitDefaultHelpCmd()
				_ = cmd.Help()
			}
		},
	}
}

func addHelpFlag(name string, fs *pflag.FlagSet) {
	fs.BoolP(flagHelp, flagHelpShortHand, false, fmt.Sprintf("Help for %s.", name))
}

func addHelpCommandFlag(usage string, fs *pflag.FlagSet) {
	fs.BoolP(
		flagHelp,
		flagHelpShortHand,
		false,
		fmt.Sprintf("Help for the %s command.", color.GreenString(strings.Split(usage, " ")[0])),
	)

}
