package main

import (
	"fmt"
	"os"

	"github.com/sdhornet/gator/internal/config"
)

const username = "nate"

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %s\n", err)
		os.Exit(1)
	}
	cfg.Print()

	if err := cfg.SetUser(username); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting username: %s", err)
		os.Exit(1)
	}
	cfg.Print()
}
