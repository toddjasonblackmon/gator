package main

import (
	"errors"
	"fmt"
	"github.com/toddjasonblackmon/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if ok {
		return handler(s, cmd)
	} else {
		return fmt.Errorf("command %s not found", cmd.name)
	}
}

func NewCommands() *commands {
	var c commands
	c.handlers = make(map[string]func(*state, command) error)
	return &c
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("invalid number of arguments given")
	}

	s.config.SetUser(cmd.args[0])
	s.config.CurrentUserName = cmd.args[0]
	fmt.Printf("Username has been set to \"%s\"\n", cmd.args[0])

	return nil
}
