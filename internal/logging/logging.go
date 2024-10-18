package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger() *log.Logger {

	logger := log.New()

	file, err := os.OpenFile("eolctl.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}

	// Set logrus to use multiWriter as the output
	logger.SetOutput(file)
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.JSONFormatter{})

	return logger
}
