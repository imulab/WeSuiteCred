package show

import "github.com/urfave/cli/v2"

type config struct {
	Query string
}

func (c *config) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "query",
			Aliases:     []string{"q"},
			Usage:       "Query that would match corporation id or name",
			Destination: &c.Query,
			Required:    true,
		},
	}
}
