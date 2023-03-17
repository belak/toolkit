package internal

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsATTY() bool {
	return isatty.IsTerminal(os.Stdin.Fd())
}
