# Gator - RSS feed aggregator [^1]

<!--toc:start-->
- [Gator - RSS feed aggregator [^1]](#gator---rss-feed-aggregator-1)
  - [Installation](#installation)
    - [Pre-requisites](#pre-requisites)
    - [Setup guide](#setup-guide)
  - [Using Gator](#using-gator)
<!--toc:end-->

Gator is a CLI utility tool that allows you to subscribe to RSS feeds
and store new entries in a database for later viewing.
It supports multiple users, each with their own feed subscriptions.
Of course, all data is stored locally in Postgres so it's persistent.

## Installation

### Pre-requisites

- Postgres (create a user and start the database service)
- Go compiler

### Setup guide

1. Clone repo
2. Navigate to the root of the project
3. Run: `go install`
4. Create a config file at: `~/.gatorconfig.json` with the following contents:

```JSON
{
 "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
 "current_user_name": "username"
}
```

Now you can run `gator <cmd> [args]` from anywhere.

## Using Gator

A simple example workflow would be:

1. Register a new user: `gator register <username>` (this also logs you in)
2. Add a new feed: `gator addfeed <name> <url>` (this also follows the feed)
3. Run `gator agg 5s` command and you will see a maximum of 20 fresh headlines
getting pulled from each followed RSS feed every 5s.
(do not run it with too small of an interval or you risk DoS-ing the server)
4. Run `gator browse [n]` to see the latest 'n' posts from your followed feeds.
In some terminals you can also shift-click the link to open it in a browser.

Below is a list of available commands with a short description.
This same list is provided if you run `gator help`

| Command | Description |
|----------|-------------|
| `register <username>`    | Register a new user|
| `login <username>`      | Log in as an existing user|
| `users`       | List all users|
| `reset`       | Delete all users and their stored data|
| `agg <time_between_reqs>`        | Run the aggregator|
| `addfeed <name> <link>`    | Add a new feed|
| `feeds`       | List all feeds|
| `follow <feed_name>`     | Follow an existing feed|
| `following`   | List feeds followed by current user|
| `unfollow <feed_name>`   | Unfollow a feed|
| `browse`      | Browse posts|
| `help`        | list available commands|

[^1]: Disclaimer: This is a [Boot.dev](Boot.dev) project
