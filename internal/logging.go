package helpers

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger() *log.Logger {

	logger := log.New()

	// Set logrus to use multiWriter as the output
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger
}
