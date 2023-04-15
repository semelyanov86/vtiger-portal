package logger

import (
	"github.com/sirupsen/logrus"
)

func Debug(msg LogMessage) {
	logrus.Debug(msg)
}

func Debugf(format string, msg LogMessage) {
	logrus.Debugf(format, msg)
}

func Info(msg LogMessage) {
	logrus.Info(msg)
}

func Infof(format string, msg LogMessage) {
	logrus.Infof(format, msg)
}

func Warn(msg LogMessage) {
	logrus.Warn(msg)
}

func Warnf(format string, msg LogMessage) {
	logrus.Warnf(format, msg)
}

func Error(msg LogMessage) {
	logrus.Error(msg)
}

func Errorf(format string, msg LogMessage) {
	logrus.Errorf(format, msg)
}
