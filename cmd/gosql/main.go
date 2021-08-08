// gosql is a tool to generate SQL queries and Go methods that execute those SQL queries and scan their results.
package main

import (
	"fmt"
	"os"

	"github.com/frk/gosql/internal/command"
	"github.com/frk/gosql/internal/config"
)

func main() {
	cfg := config.DefaultConfig
	cfg.ParseFlags(command.PrintUsage)
	if err := cfg.ParseFile(); err != nil {
		fmt.Fprintf(os.Stderr, "gosql: failed parsing config file ...\n - %v\n", err)
		os.Exit(2)
	}

	cmd, err := command.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gosql: failed to initialize the command ...\n - %v\n", err)
		os.Exit(2)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "gosql: an error occurred ...\n - %v\n", err)
		os.Exit(2)
	}
}
