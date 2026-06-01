package main

import (
	"fmt"

	"github.com/dsledge1/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	cfg.SetUser("david")
	config.Read()
}

func login() {

}

func registerUser() {

}

func users() {

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login expects a single argument: the username")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User set to: ", s.config.User)
	return nil
}

type state struct {
	config *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	login        command
	registerUser command
	users        command
	handlerMap   map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := commands.handlerMap[cmd]; ok {

	}
}

func (c *commands) register(name string, f func(*state, command) error) {

}
