package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

func Log(c context.Context) *logrus.Entry {
  return logrus.WithFields(logrus.Fields{
    "traceid": c.Value("traceid"),
  })
}
