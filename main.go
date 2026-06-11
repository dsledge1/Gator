package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dsledge1/Gator/internal/config"
	"github.com/dsledge1/Gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	s := state{config: &cfg}
	c := commands{handlerMap: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	err = c.run(&s, command{name: args[1], args: args[2:]})
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", s.config.DB_URL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	dbQueries := database.New(db)
	s.db = dbQueries

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("beep")
		return fmt.Errorf("login expects a single argument: the username")
	}
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User set to: ", s.config.User)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("register expects a single argument: the username")
	}
	userName := cmd.args[0]
	ctx := context.Background()
	uid := uuid.New().String() //Likely issue here
	t := time.Now()
	user := database.CreateUserParams{Name: userName,
		ID:        uid,
		CreatedAt: t,
		UpdatedAt: t,
	}
	u, err := s.db.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("User created: %v\n", u)
	fmt.Printf("User data: %v\n", user)
	s.config.SetUser(userName)
	return nil
}

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlerMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if runCommand, ok := c.handlerMap[cmd.name]; ok {
		return runCommand(s, cmd)

	} else {
		return fmt.Errorf("Unknown command: %s", cmd.name)
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.handlerMap[name] == nil {
		c.handlerMap[name] = f
	} else {
		fmt.Printf("Handler for %s already exists\n", name)
	}
}
