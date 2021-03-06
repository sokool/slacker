package main

import (
	"fmt"
	"log"

	"os"

	"github.com/sokool/slacker/replacer"
	"github.com/sokool/slacker/server"
)

func main() {

	server.Token = osVar("BOT_TOKEN", "")
	server.Address = osVar("BOT_ADDR", "localhost:1234")
	server.Register(replacer.OutWebHook)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func Welcome(m server.Message) (string, error) {
	return fmt.Sprintf("Eloszki, %s\nYour text is: %s\n", m.UserName, m.Text), nil
}

// Read environment variable, if empty return default
func osVar(name, def string) string {
	ev := os.Getenv(name)
	if ev == "" {
		return def
	}
	return ev
}
