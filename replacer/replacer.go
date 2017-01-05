package replacer

import (
	"regexp"

	"strings"

	"log"
	"net/http"

	"encoding/json"

	"math/rand"

	"sync"

	"errors"
	"github.com/sokool/slacker/server"
)

var (
	dictionaryURL = "http://workshop.x7467.com:1080/"
	onlyWords     = regexp.MustCompile(`~[a-z A-Z]`)
)

type (
	dRequest struct {
		Word     string   `json:"word"`
		Synonyms []string `json:"synonyms"`
	}
)

func init() {

}

func OutWebHook(m server.Message) (string, error) {
	wc, err := incomingWords(m.Text)
	if err != nil {
		return "<INFO> nothing to replace", nil
	}

	var tex string = m.Text

	for src, dst := range multipleCall(wc) {
		tex = strings.Replace(tex, src, dst, 1)
	}

	return tex, nil
}

func incomingWords(s string) ([]string, error) {
	// dłuższe niz 1 znak

	sc := strings.Split(s, " ")

	return sc, nil

}

func multipleCall(sc []string) map[string]string {
	wg := sync.WaitGroup{}
	sen := map[string]string{}
	for _, w := range sc {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			dr, err := findSynonymy(word)
			if err != nil {
				sen[w] = w
				return
			}
			sen[word] = dr.Random()
		}(w)
	}

	wg.Wait()

	return sen

}

func findSynonymy(w string) (dRequest, error) {
	dr := dRequest{}

	u := dictionaryURL + w
	resp, err := http.Get(u)
	defer resp.Body.Close()

	if err != nil {
		return dr, errors.New("Can not connect")
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Word not found: " + w)
		return dr, errors.New("Word not found:" + w)
	}

	err = json.NewDecoder(resp.Body).Decode(&dr)
	if err != nil {
		return dr, err
	}

	return dr, nil

}

func (r dRequest) Random() string {

	return r.Synonyms[rand.Intn(len(r.Synonyms))]
}
