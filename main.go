package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sdhornet/gator/internal/config"
	"github.com/sdhornet/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	return f(s, cmd)
}

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
	fmt.Printf("%+v\n", user)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %s\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB error: %s\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := state{cfg: &cfg, db: dbQueries}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: gator <command> <arguments>")
		os.Exit(1)
	}

	cmd := command{name: args[0], args: args[1:]}

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Run error: %s\n", err)
		os.Exit(1)
	}
}
