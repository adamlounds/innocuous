package main

import (
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

const (
	defaultHTTPport   = ":8080"
	defaultTelnetPort = ":3333"
)

var count = 0
var numActiveSessions = 0
var numConnections = 0
var numLinesReceived = 0
var numWordsReceived = 0
var top5Words = [5]string{"lorem", "ipsum", "dolor", "sit", "amet"}
var top5Letters = [5]string{"e", "t", "a", "o", "i"}
var serverCtx context.Context

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx = context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting Innocuous server %s", "v1.0")

	startHTTPServer(defaultHTTPport)
	startTelnetServer(defaultTelnetPort)
	select {}
}

func startTelnetServer(TelnetPort string) {
	server := tcp_server.New(TelnetPort)

	server.OnNewClient(func(c *tcp_server.Client) {
		numConnections++
		numActiveSessions++
		c.Send(fmt.Sprintf("Hello. You are connection %d\n", numActiveSessions))
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		numLinesReceived++

		words := strings.Fields(message)
		numWordsReceived += len(words)
		c.Send(fmt.Sprintf("Received %d words\n", len(words)))

	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		numActiveSessions--
	})

	go server.Listen()
	lg.Log(serverCtx).Infof("Started telnet server on port %s", TelnetPort)
}

func startHTTPServer(HTTPPort string) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
		w.Write([]byte("connections: " + strconv.Itoa(numConnections) + "<br>"))
		w.Write([]byte("lines received: " + strconv.Itoa(numLinesReceived) + "<br>"))
		w.Write([]byte("words received: " + strconv.Itoa(numWordsReceived) + "<br>"))
	})

	go func() {
		err := http.ListenAndServe(HTTPPort, r)
		if err != nil {
			panic("Listen: " + err.Error())
		}
	}()
	lg.Log(serverCtx).Infof("Started web server on port %s", HTTPPort)
}
