package main

import (
	"github.com/pressly/lg"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/net/context"
	"testing"
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
		So(func() { startHTTPServer(":9099") }, ShouldNotPanic)

	})
}
