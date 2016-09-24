package main

import (
	"log"

	"fmt"

	"github.com/sokool/slacker/server"
)

type WelcomeResponse struct {
	Text string `json:"text"`
}

func main() {
	server.Register(welcome)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func welcome(m server.Message) (interface{}, error) {
	return WelcomeResponse{fmt.Sprintf("Eloszki, %s !", m.UserName)}, nil
}
