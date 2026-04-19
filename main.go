package main

import "fmt"
import "github.com/toddjasonblackmon/gator/internal/config"

func main() {
	fmt.Println("Hello World!")

	c := config.Read()

	c.SetUser("todd")

	c = config.Read()
	fmt.Printf("%s\n%s\n", c.DbURL, c.CurrentUserName)
}
