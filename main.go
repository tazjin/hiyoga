package main

import (
	"os"

	"github.com/tazjin/hiyoga/bookings"
	"github.com/tazjin/hiyoga/classes"
	"github.com/tazjin/hiyoga/classtypes"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "HiYoga CLI"
	app.Usage = "Get moving!"

	app.Commands = []cli.Command{
		classes.ListClassesCommand(),
		classtypes.ListClassTypesCommand(),
		bookings.ListBookingsCommand(),
	}

	app.Run(os.Args)
}
