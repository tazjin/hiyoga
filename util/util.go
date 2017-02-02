package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Fail(msg string) {
	fmt.Fprintf(os.Stderr, color.RedString(msg))
	os.Exit(1)
}
