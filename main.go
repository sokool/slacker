package main

import (
	"log"

	"fmt"

	"os"

	"github.com/sokool/slacker/server"
)

func main() {

	server.Token = osEnvDefault("BOT_TOKEN", "")
	server.Address = osEnvDefault("BOT_ADDR", "localhost:1234")

	server.Register(Welcome)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func Welcome(m server.Message) (string, error) {
	return fmt.Sprintf("Eloszki, %s !", m.UserName), nil
}

// Read environment variable, if empty return def
func osEnvDefault(name, def string) string {
	ev := os.Getenv(name)
	if ev == "" {
		return def
	}
	return ev
}
