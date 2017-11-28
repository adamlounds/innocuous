package main

import (
	"bytes"
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"net"
	"testing"
	"time"
)

func TestHTTPServer(t *testing.T) {

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx = context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting Innocuous server %s", "v1.0")

	Convey("Innocuous server compiles", t, func() {
		So(func() { startHTTPServer(":9098") }, ShouldNotPanic)

	})
}

func TestTelnetServer(t *testing.T) {

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	lg.RedirectStdlogOutput(logger)
	lg.DefaultLogger = logger
	serverCtx = context.Background()
	serverCtx = lg.WithLoggerContext(serverCtx, logger)
	lg.Log(serverCtx).Infof("Starting Innocuous server %s", "v1.0")

	Convey("Innocuous telnet server can be started", t, func() {
		So(func() { startTelnetServer(":9099") }, ShouldNotPanic)
	})

	// wait for server to wake up
	time.Sleep(10 * time.Millisecond)

	var conn net.Conn
	var err error
	Convey("Can connect to server", t, func() {
		conn, err = net.Dial("tcp", "localhost:9099")
		lg.Log(serverCtx).Infof("connected %#V", conn)
		So(err, ShouldEqual, nil)
	})

	Convey("..and it responds with a welcome message", t, func() {
		reply := make([]byte, 1024)
		_, err = conn.Read(reply)

		// this one was fun :-)
		reply = bytes.Trim(reply, "\x00")
		So(err, ShouldEqual, nil)
		So(string(reply), ShouldEqual, "Hello. You are connection 1\n")
	})

	conn.Write([]byte("Test message\n"))

}
