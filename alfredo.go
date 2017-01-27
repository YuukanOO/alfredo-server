package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	var configPath string

	app := cli.NewApp()
	app.Name = "alfredo"
	app.Usage = "Flexible and light home automation server"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from toml `FILE`",
			Value:       "alfredo.toml",
			Destination: &configPath,
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "Starts the alfredo server",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}
