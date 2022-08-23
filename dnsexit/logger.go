package dnsexit

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GetLogger(name string) *logrus.Entry {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true

	log := logrus.New()

	log.SetFormatter(customFormatter)

	log.Out = os.Stdout

	logger := log.WithFields(logrus.Fields{"name": name})

	return logger
}
