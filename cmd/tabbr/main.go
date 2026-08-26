package main

import (
	"fmt"
	"os"

	"github.com/waradu/tabbr/internal/cli"
)

func main() {
	err := cli.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tabbr:", err)
		os.Exit(1)
	}
}
