package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/sdhornet/gator/internal/config"
	"github.com/sdhornet/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to db: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := &state{cfg: &cfg, db: dbQueries}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: gator <command> <arguments>")
		os.Exit(1)
	}

	cmd := command{name: args[0], args: args[1:]}

	if err := cmds.run(s, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "Run error: %s\n", err)
		os.Exit(1)
	}
}
