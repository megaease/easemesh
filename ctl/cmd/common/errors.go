package common

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// ExitWithError exits with self-defined message not the one of cobra(such as usage).
func ExitWithError(err error) {
	if err != nil {
		color.New(color.FgRed).Fprint(os.Stderr, "Error: ")
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// ExitWithErrorf wraps ExitWithError with format.
func ExitWithErrorf(format string, a ...interface{}) {
	ExitWithError(fmt.Errorf(format, a...))
}

func OutputErrorf(format string, a ...interface{}) {
	color.New(color.FgRed).Fprintf(os.Stderr, "Error: ")
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func OutputError(err error) {
	if err != nil {
		color.New(color.FgRed).Fprint(os.Stderr, "Error: ")
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
