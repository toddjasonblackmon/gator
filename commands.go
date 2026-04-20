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
	fmt.Println(user)

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

	// Add a follow record for this user
	_, err = addFollow(s, url)
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("invalid number of arguments given")
	}

	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("unable to get list of feeds")
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n * %s\n * from %s\n",
			feed.Name, feed.Url, feed.Username)
	}

	return nil
}

func addFollow(s *state, url string) (database.CreateFeedFollowRow, error) {
	var record database.CreateFeedFollowRow
	ctx := context.Background()

	feedInfo, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return record, fmt.Errorf("unable to get feed by URL: %w", err)
	}

	userInfo, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return record, fmt.Errorf("unable to get user by name: %w", err)
	}

	// Create the new following record
	ff := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userInfo.ID,
		FeedID:    feedInfo.ID,
	}

	record, err = s.db.CreateFeedFollow(ctx, ff)
	if err != nil {
		return record, fmt.Errorf("unable to add follow record: %w", err)
	}

	return record, nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("invalid number of arguments given")
	}

	url := cmd.args[0]

	record, err := addFollow(s, url)
	if err != nil {
		return err
	}
	fmt.Printf("%s following %s\n", record.UserName, record.FeedName)

	return nil

}

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("invalid number of arguments given")
	}

	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("unable to get list of feeds being followed")
	}

	for _, follow := range follows {
		fmt.Printf("%s\n", follow.Feed)
	}

	return nil

}
