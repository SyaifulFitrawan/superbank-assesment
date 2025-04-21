package utils

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	SetRequest(reqID, url string)
	SetContext(context string)
	Error(error interface{}, context ...string)
	Log(info interface{}, context ...string)
	Warn(warn interface{}, context ...string)
}

type AppLogger struct {
	logger  *logrus.Logger
	context string
}

func NewLogger() Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &AppLogger{
		logger:  logger,
		context: "account-management",
	}
}

func (l *AppLogger) SetRequest(reqID, url string) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"request": map[string]string{
			"id":  reqID,
			"url": url,
		},
	}).Logger
}

func (l *AppLogger) SetContext(context string) {
	l.context = context
}

func (l *AppLogger) Error(error interface{}, context ...string) {
	if len(context) > 0 {
		l.logger.WithFields(logrus.Fields{
			"context": context[0],
			"message": error,
		}).Error()
	} else {
		l.logger.Error(error)
	}
}

func (l *AppLogger) Log(info interface{}, context ...string) {
	if len(context) > 0 {
		l.logger.WithFields(logrus.Fields{
			"context": context[0],
			"message": info,
		}).Info()
	} else {
		l.logger.Info(info)
	}
}

func (l *AppLogger) Warn(warn interface{}, context ...string) {
	if len(context) > 0 {
		l.logger.WithFields(logrus.Fields{
			"context": context[0],
			"message": warn,
		}).Warn()
	} else {
		l.logger.Warn(warn)
	}
}
