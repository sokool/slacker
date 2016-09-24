package server

import (
	"net/http"
	"os"

	"io"

	"log"

	"io/ioutil"
	"net/url"

	"regexp"

	"encoding/json"
	"fmt"
)

var (
	address string
	token   string
	server  *defaultServer
)

type (
	defaultServer struct {
		hooks map[string]WebHook
		hook  WebHook
	}

	WebHook func(Message) (interface{}, error)

	Message struct {
		Token       string `json:"token"`
		TeamID      string `json:"team_id"`
		TeamDomain  string `json:"ream_domain"`
		ChannelID   string `json:"channel_id"`
		ChannelName string `json:"channel_name"`
		Timestamp   string `json:"timestamp"`
		UserID      string `json:"user_id"`
		UserName    string `json:"user_name"`
		Text        string `json:"text"`
		TriggerWord string `json:"trigger_word"`
	}
)

func init() {
	address, _ = osEnvDefault("BOT_ADDR", "localhost:1234")
	token, _ = osEnvDefault("BOT_TOKEN", "HA1.54")
	server = &defaultServer{}

}

func (s *defaultServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// message is created based on request body
	m, err := create(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error occoured when message was created: %s", err.Error())
		return
	}

	// Validates create message, check if all data are as server expects
	if ok, errs := isValid(m); !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Message is not valid due %s\n", errs)

		return
	}

	// Check if there is any webhook agregated in server
	if s.hook == nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("No web hock attached")
		return
	}

	// Call attached web hook and check if there is an error due message processing
	re, err := s.hook(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())

		return
	}

	// Send response to the caller
	if err := json.NewEncoder(w).Encode(re); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func Register(h WebHook) {
	server.hook = h

}

func Run() error {
	return http.ListenAndServe(address, server)
}

func isValid(m Message) (bool, []error) {
	var ec []error

	if m.Token != token {
		ec = append(ec, fmt.Errorf("Wrong token"))
	}

	return len(ec) == 0, ec

}

// Read environment variable, if empty return def
func osEnvDefault(name, def string) (string, bool) {
	ev := os.Getenv(name)
	if ev == "" {
		return def, false
	}
	return ev, true
}

func create(r io.Reader) (Message, error) {
	var m Message
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return m, err
	}

	re := regexp.MustCompile(`\r?\n`)
	input := re.ReplaceAllString(string(b), "")

	q, err := url.ParseQuery(input)
	if err != nil {
		return m, err
	}

	m = Message{
		Token:       q.Get("token"),
		TeamID:      q.Get("team_id"),
		TeamDomain:  q.Get("team_domain"),
		ChannelID:   q.Get("channel_id"),
		ChannelName: q.Get("channel_name"),
		Timestamp:   q.Get("timestamp"),
		UserID:      q.Get("user_id"),
		UserName:    q.Get("user_name"),
		Text:        q.Get("text"),
		TriggerWord: q.Get("trigger_word"),
	}

	return m, nil
}
