package main

import (
	"database/sql"
	"fmt"
	"github.com/Nic-Vladimir/blog_aggregator/internal/config"
	"github.com/Nic-Vladimir/blog_aggregator/internal/database"
	_ "github.com/lib/pq"
	"os"
)

type state struct {
	db                 *database.Queries
	config             *config.Config
	registeredCommands *commands
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	cmds := commands{
		registeredCommands: make(map[string]commandInfo),
	}

	programState := state{
		config:             cfg,
		registeredCommands: &cmds,
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: gator <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	dbUrl := cfg.DBUrl
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	programState.db = dbQueries

	cmds.register("register", "Register a new user", handlerRegister)
	cmds.register("login", "Log in as an existing user", handlerLogin)
	cmds.register("users", "List all users", handlerListUsers)
	cmds.register("reset", "Delete all users and their stored data", handlerResetUsers)
	cmds.register("agg", "Run the aggregator", handlerAggregateRss)
	cmds.register("addfeed", "Add a new feed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", "List all feeds", handlerListFeeds)
	cmds.register("follow", "Follow an existing feed", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", "List feeds followed by current user", middlewareLoggedIn(handlerListFollowing))
	cmds.register("unfollow", "Unfollow a feed", middlewareLoggedIn(handlerUnfollowFeed))
	cmds.register("browse", "Browse posts", middlewareLoggedIn(handlerBrowseFeeds))
	cmds.register("help", "list available commands", handlerHelp)

	if err := cmds.Run(&programState, cmd); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
