package util

import (
	"fmt"
	"log"
	"os"
)

const (
	DEBUG = 1 + iota
	INFO
	WARNING
	ERROR
	PANIC
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
var logLevel = INFO

func InitLogger(option Option) {
	if option.LogPath != "" {
		file, err := os.OpenFile(option.LogPath, os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			LogError("%#v", err)
		} else {
			logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		}
	}
	if option.LogLevel >= DEBUG && option.LogLevel <= PANIC {
		logLevel = option.LogLevel
	}
}

func LogPanic(format string, v ...interface{}) {
	if logLevel <= PANIC {
		format = "[PANIC] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogError(format string, v ...interface{}) {
	if logLevel <= ERROR {
		format = "[ERROR] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogWarn(format string, v ...interface{}) {
	if logLevel <= WARNING {
		format = "[WARNING] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogInfo(format string, v ...interface{}) {
	if logLevel <= INFO {
		format = "[INFO] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogDebug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		format = "[DEBUG] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}
