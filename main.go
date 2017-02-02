package main

import (
	"fmt"
	"os"
	"strconv"
	"hiyoga/classes"
	"hiyoga/util"
)

const MAJORSTUEN string = "94c207f7-fdc0-4de2-8ca4-aa42e8387b60"

func main() {
	days := getDaysArg()

	if days > 7 {
		util.Fail("Can not print more than one week in advance!")
	}

	c, err := classes.ListClasses(MAJORSTUEN, days)

	if err != nil {
		util.Fail(fmt.Sprintf("Could not list classes: %v\n", err))
	}

	classes.PrettyPrintClassResponse(days, &c)
}

// Generate the date list component of the classes list URI
func getDaysArg() int {
	args := os.Args[1:]

	if len(args) == 1 {
		days, err := strconv.Atoi(args[0])

		if err != nil {
			util.Fail("Usage: hiyoga [number of days]")
		}

		return days
	}

	return 3
}
