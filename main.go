package main

import (
	"absurdlab.io/WeSuiteCred/buildinfo"
	"absurdlab.io/WeSuiteCred/cmd/listener"
	show "absurdlab.io/WeSuiteCred/cmd/query"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	log.Info().
		Str("version", buildinfo.Version).
		Str("revision", buildinfo.Revision).
		Time("compiled_at", buildinfo.CompiledAtTime()).
		Msg("WeSuiteCred binary")

	app := &cli.App{
		Name:      "WeSuiteCred",
		Usage:     "WeSuiteCred listens for credential change messages from WeTriage and manages credentials for the suite and its authorized apps.",
		Version:   buildinfo.Version,
		Compiled:  buildinfo.CompiledAtTime(),
		Copyright: "MIT",
		Authors: []*cli.Author{
			{Name: "Weinan Qiu", Email: "davidiamyou@gmail.com"},
		},
		Commands: []*cli.Command{
			listener.Command(),
			show.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app.")
	}
}
