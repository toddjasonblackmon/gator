package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/toddjasonblackmon/gator/internal/config"
	"github.com/toddjasonblackmon/gator/internal/database"
	"os"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	progState := initState()

	dbURL := progState.config.DbURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Unable to connect to the database")
		os.Exit(1)
	}

	dbQueries := database.New(db)
	progState.db = dbQueries

	commandTable := NewCommands()
	commandTable.register("login", handlerLogin)
	commandTable.register("register", handlerRegister)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	c := command{name: os.Args[1], args: os.Args[2:]}

	err = commandTable.run(&progState, c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func initState() state {
	var s state
	currentConfig := config.Read()
	s.config = &currentConfig

	return s
}
