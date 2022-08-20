package dnsexit

import (
	"os"

	"github.com/sirupsen/logrus"
)

func GetLogger(name string) *logrus.Entry {
	log := logrus.New()

	log.Out = os.Stdout

	logger := logrus.WithFields(logrus.Fields{"name": name})

	return logger
}