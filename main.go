package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sdhornet/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Missing username")
	}
	username := cmd.args[2]

	if err := s.cfg.SetUser(username); err != nil {
		return err
	}

	fmt.Printf("Username set to: %s\n", username)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %s\n", err)
		os.Exit(1)
	}
	cfg.Print()

	if err := cfg.SetUser(username); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting username: %s\n", err)
		os.Exit(1)
	}

	cfg2, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %s\n", err)
		os.Exit(1)
	}
	cfg2.Print()
}
