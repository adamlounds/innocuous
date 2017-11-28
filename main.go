package main

import (
	"encoding/json"
	"fmt"
	"github.com/firstrow/tcp_server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"strings"
)

type statistics struct {
	numLinesReceived int      `json:"numLines"`
	numWordsReceived int      `json:"count"`
	top5Words        []string `json:"top_5_words"`
	top5Letters      []string `json:"top_5_letters"`
}

const (
	defaultHTTPport   = ":8080"
	defaultTelnetPort = ":3333"
)

var numActiveSessions = 0
var numConnections = 0
var stats statistics
var allWords map[string]int
var allLetters map[string]int

//var top5Words = [5]string{"lorem", "ipsum", "dolor", "sit", "amet"}
//var top5Letters = [5]string{"e", "t", "a", "o", "i"}
var serverCtx context.Context

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx = context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting Innocuous server %s", "v1.0")
	stats = statistics{0, 0, []string{"lorem", "ipsum", "dolor", "sit", "amet"}, []string{"e", "t", "a", "o", "i"}}

	allWords = make(map[string]int)
	allLetters = make(map[string]int)

	startHTTPServer(defaultHTTPport)
	wordsChan := startTelnetServer(defaultTelnetPort)

	defer close(wordsChan)

	go func() {
		lg.Log(serverCtx).Infof("parser started")
		words := <-wordsChan
		for i := 0; i < len(words); i++ {
			word := words[i]
			allWords[word]++
			wordCount := allWords[word]
			lg.Log(serverCtx).Infof("seen %s %d times", word, wordCount)

			letters := strings.Split(word, "")
			for j := 0; j < len(letters); j++ {
				letter := letters[j]
				allLetters[letter]++
				letterCount := allLetters[letter]
				lg.Log(serverCtx).Infof("seen %s %d times", letter, letterCount)
			}

			// if count is greater than current minimum-number-to-qualify,
			// check existing top-5s to recalc minimum-number-to-qualify
			// and evict the lowest if needed
		}
	}()

	// cannot get this to work with telnet server :-(
	wordsChan <- []string{"lorem", "ipsum", "dolor", "sit", "amet", "lorem"}

	select {}
}

func startTelnetServer(TelnetPort string) chan []string {
	server := tcp_server.New(TelnetPort)
	wordsChan := make(chan []string)

	server.OnNewClient(func(c *tcp_server.Client) {
		numConnections++
		numActiveSessions++
		c.Send(fmt.Sprintf("Hello. You are connection %d\n", numActiveSessions))
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		stats.numLinesReceived++

		words := strings.Fields(message)
		stats.numWordsReceived += len(words)
		c.Send(fmt.Sprintf("Received %d words\n", len(words)))

		wordsChan <- []string{"two", "three"}

	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		numActiveSessions--
	})

	go server.Listen()
	lg.Log(serverCtx).Infof("Started telnet server on port %s", TelnetPort)

	return wordsChan
}

func startHTTPServer(HTTPPort string) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.Write([]byte("connections: " + strconv.Itoa(numConnections) + "<br>"))
		w.Write([]byte("lines received: " + strconv.Itoa(stats.numLinesReceived) + "<br>"))
		w.Write([]byte("words received: " + strconv.Itoa(stats.numWordsReceived) + "<br>"))

		// aaagh, json eludes me once more
		json, err := json.Marshal(stats)
		if err != nil {
			lg.Log(serverCtx).Errorf("cannot marshal %s", err)
		}
		lg.Log(serverCtx).Infof("json %s", string(json))
		lg.Log(serverCtx).Infof("stats %+v", stats)

	})

	go func() {
		err := http.ListenAndServe(HTTPPort, r)
		if err != nil {
			panic("Listen: " + err.Error())
		}
	}()
	lg.Log(serverCtx).Infof("Started web server on port %s", HTTPPort)
}
