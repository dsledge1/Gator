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
