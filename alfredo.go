package main

import (
	"os"

	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/handlers"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
)

func main() {
	var configPath string
	var adaptersPath string

	app := cli.NewApp()
	app.Name = "alfredo"
	app.Usage = "Flexible and light home automation server"
	app.Version = env.Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from toml `FILE`",
			Value:       "./alfredo.toml",
			Destination: &configPath,
		},
		cli.StringFlag{
			Name:        "adapters, a",
			Usage:       "Load adapters from json `FILE`",
			Value:       "./adapters.json",
			Destination: &adaptersPath,
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:  "run",
			Usage: "Starts the alfredo server",
			Action: func(c *cli.Context) error {
				if err := env.LoadFromFile(configPath); err != nil {
					return err
				}

				defer env.Cleanup()

				if err := env.LoadAdapters(adaptersPath); err != nil {
					return err
				}

				r := gin.Default()

				r.Use(middlewares.CORS(&middlewares.CORSOptions{
					AllowedOrigins: env.Current().Server.AllowedOrigins,
					AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
				}), middlewares.DB())

				handlers.Register(r)

				return r.Run(env.Current().Server.Listen)
			},
		},
		cli.Command{
			Name:  "reset",
			Usage: "Drops the database to start from a fresh start. Used mostly for dev.",
			Action: func(c *cli.Context) error {
				if err := env.LoadFromFile(configPath); err != nil {
					return err
				}

				sess, db := env.Current().Database.GetSession()
				defer sess.Close()
				return db.DropDatabase()
			},
		},
	}

	app.Run(os.Args)
}
