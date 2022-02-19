package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/paulwrubel/moneybags-server/constants"
	"github.com/paulwrubel/moneybags-server/database"
	"github.com/paulwrubel/moneybags-server/injection"
	"github.com/paulwrubel/moneybags-server/routing"
	log "github.com/sirupsen/logrus"
)

func main() {
	initLogger()

	log.Info("starting moneybags server")
	log.Debugf("number of CPUs: %d", runtime.NumCPU())

	dbInfo, err := getDBInfo()
	if err != nil {
		log.WithError(err).Fatal("error getting db info")
	}
	db, err := database.InitDB(dbInfo)
	if err != nil {
		log.WithError(err).Fatal("error initializing db connection")
	}
	injector := &injection.PostgresInjector{
		Database: db,
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

func getDBInfo() (*database.DBInfo, error) {
	log.Info("getting DB info")
	pgHost, isSet := os.LookupEnv(constants.PostgresHostnameEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresHostnameEnvironmentKey)
	}
	pgUser, isSet := os.LookupEnv(constants.PostgresUsernameEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresUsernameEnvironmentKey)
	}
	pgPass, isSet := os.LookupEnv(constants.PostgresPasswordEnvironmentKey)
	if !isSet {
		return nil, fmt.Errorf("environment variable %s not set", constants.PostgresPasswordEnvironmentKey)
	}

	return &database.DBInfo{
		Host:     pgHost,
		Username: pgUser,
		Password: pgPass,
	}, nil
}
