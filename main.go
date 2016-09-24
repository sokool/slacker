package main

import (
	"log"

	"fmt"

	"github.com/sokool/slacker/server"
)

func main() {
	server.Register(welcome)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func welcome(m server.Message) (string, error) {
	return fmt.Sprintf("Eloszki, %s !", m.UserName), nil
}
