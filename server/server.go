package server

import (
	"net/http"

	"io"

	"log"

	"io/ioutil"
	"net/url"

	"regexp"

	"encoding/json"
	"errors"
)

var (
	Address string = "localhost:1234"
	Token   string = ""
	srv     *server
)

type (
	WebHook func(Message) (string, error)

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

	response map[string]interface{}

	server struct {
		hook WebHook
	}
)

func init() {
	srv = &server{}
}

func isValid(m Message) (bool, []error) {
	var ec []error
	if Token != "" && m.Token != Token {
		ec = append(ec, errors.New("Given token is not valid"))
	}

	return len(ec) == 0, ec
}

//
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

func Register(h WebHook) {
	srv.hook = h
}

func Run() error {

	log.Printf("Running server on %s\n", Address)
	log.Printf("Listen for slack messages with %s token\n", Token)

	return http.ListenAndServe(Address, srv)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	o, err := s.hook(m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())

		return
	}

	log.Printf("Received message from user: %s, channel: %s, triggers: %s\n", m.UserName, m.ChannelName, m.TriggerWord)
	re := response{"text": o}
	// Send response to the caller
	if err := json.NewEncoder(w).Encode(re); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("Response served: %v\n", re["text"])

}
