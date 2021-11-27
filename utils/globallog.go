package utils

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	cyan   = 36
	gray   = 37
)

var GlobalLog *logrus.Logger

type customFormatter struct{}

func (cf *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = new(bytes.Buffer)
	}

	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	case logrus.InfoLevel:
		levelColor = cyan
	default:
		levelColor = cyan
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	if entry.HasCaller() {
		fileName := filepath.Base(entry.Caller.File)
		functionName := filepath.Base(entry.Caller.Function)
		newLog = fmt.Sprintf("\u001B[%dm[%s]\u001B[0m\u001B[%dm[%s]\u001B[0m\u001B[%dm[%s:%d %s]\u001B[0m \u001B[%dm%s\u001B[0m\n", green, timestamp, levelColor, entry.Level, blue, fileName, entry.Caller.Line, functionName, levelColor, entry.Message)
	} else {
		newLog = fmt.Sprintf("\u001B[%dm[%s]\u001B[0m\u001B[%dm[%s]\u001B[0m \u001B[%dm%s\u001B[0m\n", green, timestamp, levelColor, entry.Level, levelColor, entry.Message)
	}
	b.WriteString(newLog)
	return b.Bytes(), nil
}
