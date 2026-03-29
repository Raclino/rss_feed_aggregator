# gator

gator is a CLI RSS aggregator written in Go. It lets you register users, add and follow feeds, fetch posts from RSS feeds, and browse recent posts from the terminal.

## Requirements

To run gator, you need:

- Go
- PostgreSQL

You can either:

- use a local PostgreSQL instance
- or start PostgreSQL with the provided `docker-compose.yml`

## Install

Once the repository is on GitHub, install the CLI with:

```bash
go install github.com/your-github-username/gator@latest
```

Then run it with:

```bash
gator
```

`go run .` is only for development. After go install or go build, use the compiled gator binary.

## Config

gator reads its config from:

```bash
~/.gatorconfig.json
```

Example with a local PostgreSQL running on port 5432:

```bash
{
  "db_url": "postgres://postgres:superpswd@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

If you use Docker Compose and expose PostgreSQL on port 5433, use:

```bash
{
  "db_url": "postgres://postgres:superpswd@localhost:5433/gator?sslmode=disable",
  "current_user_name": ""
}
```

## Database setup

If PostgreSQL is running locally, create a gator database and run:

```bash
make mig-up
```

If you use Docker Compose:

```bash
docker-compose up -d
make mig-up
```

Make sure your `db_url` matches the port, user, password, and database name from your PostgreSQL setup.

## Usage

Examples:

```bash
gator register Thomas
gator login Thomas
gator addfeed hackerNews https://hnrss.org/newest
gator follow https://hnrss.org/newest
gator following
gator agg 30s
gator browse
gator browse 10
```

## Notes

agg runs in a loop and fetches feeds on the interval you provide
browse shows recent posts for the current user
`go run .` is for local development
gator is the production CLI binary
