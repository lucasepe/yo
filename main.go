package main

import (
	"fmt"
	"os"

	"github.com/lucasepe/yo/cmd"
)

func main() {
	app := cmd.Run()
	if err := app.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "mkobj error: %s\n", err.Error())
		os.Exit(1)
	}
}
