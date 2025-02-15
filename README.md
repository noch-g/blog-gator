# blog-gator

## What is this?

**blog-gator** is an implementation of the guided project [Build a Blog Aggregator](https://www.boot.dev/courses/build-blog-aggregator) from the [boot.dev](https://www.boot.dev) platform.

It is a CLI tool that scrapes RSS feeds and saves the posts to a postgresql database, allowing users to browse the posts in a terminal.


## Requirements

- Go
  - Documentation to install Go can be found [here](https://go.dev/doc/install)
- PostgreSQL
  - Documentation to install PostgreSQL can be found [here](https://www.postgresql.org/download/)

## Installation

```bash
> go install github.com/noch-g/blog-gator@latest
```

Make sure you create a `gator` database on your PostgreSQL server:
```bash
> createdb gator
```

## Configuration
The application needs to be configured with a database connection string. This can be done by creating a `.gatorconfig.json` file in the `$HOME` directory.


The config file should look like this:

```
{
    "db_url": "postgres://<user>:@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
}
```
Replace `<user>` with the user you want to use to connect to the database.
Replace `gator` with the name of the database if you chose to name it differently.



## Usage

- **Register** a user
```bash
> blog-gator register <username>
```

- **Add** feeds
```bash
> blog-gator addfeed <name> <url>
```
for example:
```bash
> blog-gator addfeed boot.dev https://blog.boot.dev/index.xml
```

- **Follow** a feed
```bash
> blog-gator follow <url>
```

- **Aggregate** the data
```bash
> blog-gator agg 1m
```
This will rotate through the existing feeds and every minute fetch the new posts from the feed that hasn't been fetched the longest.

- **Browse** the posts
```bash
> blog-gator browse [<limit>]
```
This will fetch the latest posts from the database from the user's followed feeds and display them in the terminal. If a limit is provided, it will fetch the latest posts up to that limit.
