package main

import (
	"fmt"

	"github.com/noch-g/blog-gator/internal/config"
)

func main() {
	fmt.Println("Hello, World!")

	config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println(config.CurrentUserName)
	fmt.Println(config.DbUrl)
}
