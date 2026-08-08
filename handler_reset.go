package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("reset takes zero arguments")
	}
	if err := s.db.Reset(context.Background()); err != nil {
		return fmt.Errorf("reset failed: %w", err)
	}
	fmt.Println("reset successful")
	return nil
}
