package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/toddjasonblackmon/gator/internal/database"
	"time"
)

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

	username := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("unknown user %s", username)
	}

	s.config.SetUser(username)
	s.config.CurrentUserName = username
	fmt.Printf("Username has been set to \"%s\"\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("invalid number of arguments given")
	}
	ctx := context.Background()

	username := cmd.args[0]

	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		return errors.New("attempting to create duplicate user")
	}

	user, err := s.db.CreateUser(ctx,
		database.CreateUserParams{
			uuid.New(), time.Now(),
			time.Now(), username})
	if err != nil {
		return err
	}

	s.config.SetUser(username)
	fmt.Print(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("invalid number of arguments given")
	}

	ctx := context.Background()

	return s.db.DeleteUsers(ctx)

}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("invalid number of arguments given")
	}

	ctx := context.Background()
	// s.config.CurrentUserName

	users, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("unable to get list of users")
	}

	for _, user := range users {
		selection := ""
		if user.Name == s.config.CurrentUserName {
			selection = " (current)"
		}
		fmt.Printf("* %s%s\n", user.Name, selection)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
    url := "https://www.wagslane.dev/index.xml"

	if len(cmd.args) != 0 {
		return errors.New("invalid number of arguments given")
	}
	ctx := context.Background()

    feed, err := fetchFeed(ctx, url)
    if err != nil {
        return err
    }

    fmt.Println(*feed)

    return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("invalid number of arguments given")
	}
	ctx := context.Background()

	currentUser, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feedname := cmd.args[0]
    url := cmd.args[1]




	feed, err := s.db.CreateFeed(ctx,
		database.CreateFeedParams{
			uuid.New(), time.Now(), time.Now(), 
            feedname, url, currentUser.ID})
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}


