package logger

import (
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// InitLogger ...
func InitLogger(logDir, logFile string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	os.Mkdir(logDir, os.ModeAppend)
	filepath := path.Join(logDir, logFile)
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
}

// LogError log it
func LogError(action, file, data string, err error) {
	logrus.WithFields(
		logrus.Fields{
			"action": action,
			"file":   file,
			"data":   data,
		},
	).Error(err)
}
