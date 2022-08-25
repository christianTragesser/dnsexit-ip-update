package dnsexit

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = GetLogger()
var cliLogFields = logrus.Fields{}
var updateLogFields = logrus.Fields{"component": "update"}
var recordLogFields = logrus.Fields{"component": "record"}

func GetLogger() *logrus.Logger {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true

	log := logrus.New()

	log.SetFormatter(customFormatter)

	log.Out = os.Stdout

	return log
}
