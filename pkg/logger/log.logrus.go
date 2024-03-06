package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	// create a new logger instance
	log = logrus.New()

	// set output to stdout
	log.SetOutput(os.Stdout)

	// set log formatter
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func Info(message string, fields logrus.Fields) {
	log.WithFields(fields).Info(message)
}

func Debug(message string, fields logrus.Fields) {
	log.WithFields(fields).Debug(message)
}

func Error(message string, fields logrus.Fields) {
	log.WithFields(fields).Error(message)
}

func Fatal(message string, fields logrus.Fields) {
	log.WithFields(fields).Fatal(message)
}

func Panic(message string, fields logrus.Fields) {
	log.WithFields(fields).Panic(message)
}
