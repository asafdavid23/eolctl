package helpers

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger() *log.Logger {

	logger := log.New()

	// Set logrus to use multiWriter as the output
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.JSONFormatter{})

	return logger
}
