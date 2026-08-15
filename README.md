# gator

A CLI RSS feed aggregator written in Go. gator lets you register feeds, follow the ones you care about, continuously scrape them in the background, and browse the collected posts from your terminal.

Built as part of the [boot.dev](https://www.boot.dev) backend track.

## Requirements

- **Go** 1.26 or later — needed to install the CLI (the compiled binary runs without it)
- **PostgreSQL** — the storage backend (developed against Postgres 17)
- **goose** — runs the database migrations: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Install

```bash
go install github.com/sdhornet/gator@latest
```

This compiles and installs the `gator` binary into your `$GOPATH/bin` (usually `~/go/bin`). Make sure that directory is on your `PATH`.

## Setup

**1. Create the database:**

```bash
createdb gator
```

**2. Run the migrations** from the repo root, pointing goose at your database:

```bash
goose -dir sql/schema postgres "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable" up
```

**3. Create the config file** at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://<user>:<password>@localhost:5432/gator?sslmode=disable"
}
```

`current_user_name` is managed by the program — `register` and `login` write it for you.

## Usage

```bash
gator register <name>        # create a user and log in as them
gator login <name>           # switch to an existing user
gator users                  # list users, marking the current one
gator addfeed <name> <url>   # add a feed and follow it
gator feeds                  # list all feeds and who added them
gator follow <url>           # follow an existing feed
gator following              # list the feeds you follow
gator unfollow <url>         # unfollow a feed
gator agg <interval>         # run the aggregator loop, e.g. gator agg 1m
gator browse [limit]         # show the newest posts from feeds you follow (default 2)
gator reset                  # wipe all users (and everything they own) — dev convenience
```

The intended workflow uses two terminals: leave `gator agg 1m` running in one (it visits your feeds in rotation, saving new posts and skipping duplicates), and use `gator browse` in the other to read what it has collected. Be polite to feed servers — dont run `agg` with very short intervals.
