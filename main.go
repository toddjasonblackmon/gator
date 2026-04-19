package main

import (
	"fmt"
	"github.com/toddjasonblackmon/gator/internal/config"
	"os"
)

func main() {
	progState := initState()

	commandTable := NewCommands()
	commandTable.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	c := command{name: os.Args[1], args: os.Args[2:]}

	err := commandTable.run(&progState, c)
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
