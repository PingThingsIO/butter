package main

import (
	"fmt"
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
			Usage: "List collections for a BTrDB endpoint. If only one collection is returned, its streams will be listed.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "prefix,p",
					Usage: "Prefix to filter collection names",
				},
			},
			Action: func(c *cli.Context) error {
				// Hack for positional arguments https://github.com/urfave/cli/pull/140#issuecomment-131841364
				cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", "[endpoint (default: localhost:4410)]", -1)
				endpoint := c.Args().First()
				if endpoint == "" {
					endpoint = "localhost:4410"
				}

				list(endpoint, c.String("prefix"))
				return nil
			},
		},
		{
			Name:  "rm",
			Usage: "Remove a stream from BTrDB",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes,y",
					Usage: "Skips confirmation prompt",
				},
			},
			Action: func(c *cli.Context) error {
				cli.CommandHelpTemplate = strings.Replace(cli.CommandHelpTemplate, "[arguments...]", "[endpoint] [stream uuid]", -1)
				endpoint := c.Args().First()
				if endpoint == "" {
					fmt.Println("endpoint positional argument is required")
					cli.ShowCommandHelp(c, "rm")
					os.Exit(1)
				}

				uuid := c.Args().Get(1)
				if uuid == "" {
					fmt.Println("stream uuid positional argument is required")
					cli.ShowCommandHelp(c, "rm")
					os.Exit(1)
				}

				remove(endpoint, uuid, c.Bool("yes"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
