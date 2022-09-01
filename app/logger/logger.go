package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

func WithContext(c context.Context) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"traceid": c.Value("traceid"),
	})
}
