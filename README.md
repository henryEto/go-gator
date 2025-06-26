# go-gator
An RSS feed aggregator in Go!  🐊

'go-gator' is an RSS aggregator built with Go. It allows users to follow RSS feeds, browse posts, and manage their subscriptions directly from the command line.

## Prerequisites

To run 'go-gator', you'll need the following installed on your system:

* **Go**: 'go-gator' is written in Go, so you'll need a Go development environment set up. You can download and install Go from the official Go website: [https://golang.org/doc/install](https://golang.org/doc/install)
* **PostgreSQL**: 'go-gator' uses PostgreSQL as its database. You'll need a PostgreSQL server running and accessible. You can download PostgreSQL from its official website: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)

## Installation

You can install the 'gator' CLI tool using `go install`:

```bash
go install [github.com/henryEto/go-gator@latest](https://github.com/henryEto/go-gator@latest)
````

This command will download the 'go-gator' source code, compile it, and place the `go-gator` executable in your `$GOPATH/bin` directory (or `$GOBIN` if set), making it available in your system's PATH.

## Setup and Running

1.  **Database Configuration**:
    'go-gator' requires a `config.json` file to connect to your PostgreSQL database. This file should be placed in the same directory where you plan to run the `go-gator` command. Create a `config.json` file with the following structure:

    ```json
    {
      "db_url": "your_postgresql_connection_string",
      "username": ""
    }
    ```

    Replace `"your_postgresql_connection_string"` with your actual PostgreSQL connection string. For example: `postgres://user:password@host:port/database_name?sslmode=disable`.

2.  **Initialize the Database**:
    Before you can use 'go-gator', you need to initialize its database schema. You can do this by running the `reset` command. This command will create the necessary tables in your PostgreSQL database.

    ```bash
    go-gator reset
    ```

3.  **Run 'go-gator'**:
    Once the database is set up, you can run 'go-gator' commands.

    ```bash
    go-gator <command> [args...]
    ```

## Commands

Here are a few essential commands you can run with 'go-gator':

  * **`register <username>`**: Creates a new user and sets them as the current active user.

    ```bash
    go-gator register alice
    ```

  * **`login <username>`**: Switches the current active user to an existing user.

    ```bash
    go-gator login alice
    ```

  * **`addfeed <feed_name> <url>`**: Adds a new RSS feed to the system and automatically follows it for the current user.

    ```bash
    go-gator addfeed "Go Blog" "[https://blog.golang.org/feed.atom](https://blog.golang.org/feed.atom)"
    ```

  * **`following`**: Lists all the RSS feeds the current user is following.

    ```bash
    go-gator following
    ```

  * **`browse [limit]`**: Displays the latest posts from all followed feeds. You can optionally specify a limit for the number of posts to display.

    ```bash
    go-gator browse 5
    ```

  * **`agg <time_between_reqs>`**: Starts an aggregator that periodically fetches new posts from all feeds. `time_between_reqs` can be specified in formats like `1s`, `1m`, `1h`.

    ```bash
    go-gator agg 30s
    ```

  * **`unfollow <url>`**: Unfollows a specific feed by its URL.

    ```bash
    go-gator unfollow [https://blog.golang.org/feed.atom](https://blog.golang.org/feed.atom)
    ```

  * **`feeds`**: Lists all feeds currently registered in the system (not just those followed by the current user).

    ```bash
    go-gator feeds
    ```

  * **`users`**: Lists all registered users.

    ```bash
    go-gator users
    ```

Feel free to explore other commands and functionalities\!

```
