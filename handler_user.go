package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdhornet/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	if len(cmd.args) > 1 {
		return errors.New("login takes only one username")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("Username set to: %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing username")
	}

	if len(cmd.args) > 1 {
		return errors.New("register takes only one username")
	}
	now := time.Now()
	p := database.CreateUserParams{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, Name: cmd.args[0]}
	user, err := s.db.CreateUser(context.Background(), p)
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Printf("%s was created\n", user.Name)
	fmt.Printf("%s	%s\n", user.ID, user.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) > 1 {
		return errors.New("users takes zero arguments")
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) == 0 {
		fmt.Println("no users")
		return nil
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}

	return nil
}
