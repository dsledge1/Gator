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
	cfg, err := config.Read() //read config
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}
	s := state{config: &cfg}

	//open DB
	db, err := sql.Open("postgres", s.config.DB_URL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer db.Close()
	//Create DB Queries
	dbQueries := database.New(db)
	//Create state with config and db
	s.db = dbQueries
	//Register Commands
	c := commands{handlerMap: make(map[string]func(*state, command) error)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	//Parse Args
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	//run command
	err = c.run(&s, command{name: args[1], args: args[2:]})
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("beep")
		return fmt.Errorf("login expects a single argument: the username")
	}
	ctx := context.Background()
	loginName := cmd.args[0]
	//if not in db, error and exit
	if _, err := s.db.GetUser(ctx, loginName); err != nil {
		return fmt.Errorf("User not found")
	}
	err := s.config.SetUser(loginName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	uid := uuid.New() //Likely issue here
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

func (c *commands) reset(s *state) error {
	err := s.db.ClearTable(context.Background())
	if err != nil {
		return fmt.Errorf("Error clearing table: %v", err)
	}
	return nil
}
