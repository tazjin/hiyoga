package main

import (
	"os"

	"github.com/tazjin/hiyoga/classes"
	"github.com/tazjin/hiyoga/classtypes"
	"github.com/tazjin/hiyoga/util"
	"github.com/urfave/cli"
)

// Right now HiYoga Majorstuen is the only center
const MAJORSTUEN string = "94c207f7-fdc0-4de2-8ca4-aa42e8387b60"

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "HiYoga CLI"
	app.Usage = "Get moving!"

	app.Commands = []cli.Command{
		{
			Name:    "list-classes",
			Usage:   "list upcoming yoga classes",
			Aliases: []string{"lc"},
			Action: func(c *cli.Context) error {
				days := c.Int("days")
				listAndPrintClasses(MAJORSTUEN, days)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "days, d",
					Usage: "number of days to list (including today)",
					Value: 3,
				},
			},
		},
		{
			Name:    "list-class-types",
			Usage:   "list available yoga class types",
			Aliases: []string{"lct"},
			Action: func(c *cli.Context) error {
				listAndPrintClassTypes()
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func listAndPrintClasses(center string, days int) {
	c, err := classes.ListClasses(MAJORSTUEN, days)

	if err != nil {
		util.Fail(err)
	}

	classes.PrettyPrintClassResponse(days, &c)
}

func listAndPrintClassTypes() {
	ct, err := classtypes.ListClassTypes()

	if err != nil {
		util.Fail(err)
	}

	classtypes.PrettyPrintClassTypeResponse(&ct)
}
