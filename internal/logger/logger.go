package logger

import (
	"fmt"
	"io"
	"library-api-user/internal/config"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Debug(message string, fields map[string]interface{})
}

type LoggerImpl struct {
	*logrus.Logger
	LogFile string
}

func NewLogger(logFile string) (Logger, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	env := config.ENV.Environtment

	if env == "development" {
		// Logs to console in development mode
		log.SetOutput(os.Stdout)
		log.SetLevel(logrus.DebugLevel)
	} else {
		// Write logs to a file in production mode
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		multiWriter := io.MultiWriter(file, os.Stdout)
		log.SetOutput(multiWriter)
		log.SetLevel(logrus.InfoLevel)
	}

	return &LoggerImpl{
		Logger:  log,
		LogFile: logFile,
	}, nil
}

func (l *LoggerImpl) Info(message string, fields map[string]interface{}) {
	l.WithFields(logrus.Fields(fields)).Info(message)
}

func (l *LoggerImpl) Error(message string, fields map[string]interface{}) {
	l.WithFields(logrus.Fields(fields)).Error(message)
}

func (l *LoggerImpl) Warn(message string, fields map[string]interface{}) {
	l.WithFields(logrus.Fields(fields)).Warn(message)
}

func (l *LoggerImpl) Debug(message string, fields map[string]interface{}) {
	l.WithFields(logrus.Fields(fields)).Debug(message)
}
