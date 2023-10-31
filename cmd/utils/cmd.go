package utils

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:    "utils",
		Aliases: []string{"util"},
		Usage:   "A list of utilities",
		Subcommands: []*cli.Command{
			simulateChangeAuth(),
		},
	}
}
