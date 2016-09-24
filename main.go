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

	server.Register(Welcome)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func Welcome(m server.Message) interface{} {
	return WelcomeResponse{fmt.Sprintf("Eloszki, %s !", m.UserName)}
}
