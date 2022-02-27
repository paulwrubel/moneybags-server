package main

import (
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/paulwrubel/moneybags-server/config"
	"github.com/paulwrubel/moneybags-server/injection"
	"github.com/paulwrubel/moneybags-server/routing"
	log "github.com/sirupsen/logrus"
)

func main() {
	initLogger()

	log.Info("starting moneybags server")
	log.Debugf("number of CPUs: %d", runtime.NumCPU())

	appInfo, err := config.InitializeApp()
	if err != nil {
		log.WithError(err).Fatal("error initializing app")
	}
	injector := &injection.Injector{
		AppInfo: appInfo,
	}

	log.Info("starting API server")
	routing.RunServer(injector)

	log.Info("blocking until signalled to shutdown")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan

	log.Info("shutting down")
	os.Exit(0)
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	switch strings.ToUpper(os.Getenv("MONEYBAGS_LOG_LEVEL")) {
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}
