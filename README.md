# Gator

A Go CLI for collecting RSS posts in PostgreSQL and browsing posts from your subscribed feeds.

## Requirements

- [Go](https://go.dev/doc/install) (see `go.mod` for the required version).
- [PostgreSQL](https://www.postgresql.org/download/), installed and running.
- [Goose](https://github.com/pressly/goose) for database migrations.

## Installation

Clone this repository, open its root directory, and install the CLI:

```bash
go install .
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Ensure your Go binary directory (`GOBIN`, or `$(go env GOPATH)/bin` by default) is on your `PATH`.

## Configuration and database

The examples below assume the config file is `~/.gatorconfig.json` and migrations are in `sql/schema`. Adjust these paths if your project uses different locations.

Connect as a PostgreSQL administrator:

```bash
psql -U postgres
```

Create a local application role and database, then exit:

```sql
CREATE USER gator WITH PASSWORD 'change-me';
CREATE DATABASE gator OWNER gator;
\q
```

Create `~/.gatorconfig.json` using your database credentials:

```json
{
  "db_url": "postgres://gator:change-me@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Apply migrations from the repository root using the same connection URL:

```bash
goose -dir sql/schema postgres "postgres://gator:change-me@localhost:5432/gator?sslmode=disable" up
```

## Usage

Register a user and add a feed (adding a feed also follows it):

```bash
gator register alex
gator addfeed "Hacker News" "https://hnrss.org/newest"
```

Run the scraper to collect posts at a regular interval:

```bash
gator agg 1m
```

Keep it running and browse saved posts in another terminal:

```bash
gator browse       # Show the latest 2 posts from followed feeds
gator browse 10    # Show up to 10 posts
gator login alex   # Switch the current user
```

Stop the scraper with `Ctrl+C`.
