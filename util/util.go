package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/polydawn/meep"
)

type HiyogaError struct {
	meep.AllTraits
}

func Fail(cause error) {
	hiyogaError := meep.New(
		&HiyogaError{},
		meep.Cause(cause),
	)
	fmt.Fprintf(os.Stderr, color.RedString("%v\n", hiyogaError))
	os.Exit(1)
}
