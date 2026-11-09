package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := buildConfig(os.Args)
	if err != nil {
		fmt.Printf("Could not initialize config: %s\n", err)
		os.Exit(1)
	}

	cfg.cleanup()
	fmt.Println("Thank you for playing, come again.")
}
