package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func NewLogger(logLevel string) *log.Logger {

	logger := log.New()

	level, err := log.ParseLevel(strings.ToLower(logLevel))

	if err != nil {
		logger.SetLevel(log.InfoLevel)
		logger.Warnf("Invalid log level '%s', defaulting to Info level.", logLevel)
	} else {
		logger.SetLevel(level)
	}

	// Set logrus to use multiWriter as the output
	logger.SetOutput(os.Stderr)
	// Set the log formatter with timestamp and log level
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,                  // Enables full timestamp
		TimestampFormat: "2006-01-02 15:04:05", // Custom timestamp format
	})

	return logger
}
