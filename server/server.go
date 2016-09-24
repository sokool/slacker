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

	WebHook func(Message) interface{}

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
	server = &defaultServer{
	//hooks: map[string]WebHook{},
	}

}

func (s *defaultServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m, err := parse(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if !isTokenValid(m) {
		log.Printf("Wrong token %s\n", m.Token)
		return
	}

	re := s.hook(m)
	// Send personalised response
	json.NewEncoder(w).Encode(re)

}

func Register(h WebHook) {
	server.hook = h

}

func Run() error {
	return http.ListenAndServe(address, server)
}

func isTokenValid(m Message) bool {
	return m.Token == token
}

// Read environment variable, if empty return def
func osEnvDefault(name, def string) (string, bool) {
	ev := os.Getenv(name)
	if ev == "" {
		return def, false
	}
	return ev, true
}

func parse(r io.Reader) (Message, error) {
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
