package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "butter"
	app.Usage = "Useful BTrDB CLI tools for development"

	app.Commands = []cli.Command{
		{
			Name:  "ls",
			Usage: "List collections and streams for a BTrDB endpoint",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "prefix,p"},
			},
			Action: func(c *cli.Context) error {
				cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", "[endpoint (default: localhost:4410)]", -1)
				endpoint := c.Args().First()
				if endpoint == "" {
					endpoint = "localhost:4410"
				}

				list(endpoint, c.String("prefix"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
