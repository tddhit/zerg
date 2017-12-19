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
	FATAL
	PANIC
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
var logLevel = INFO

func InitLogger(option Option) {
	if option.LogPath != "" {
		file, err := os.OpenFile(option.LogPath, os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			LogError("failed open file: %s, %s", option.LogPath, err)
		} else {
			logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
		}
	}
	if option.LogLevel >= DEBUG && option.LogLevel <= PANIC {
		logLevel = option.LogLevel
	}
}

func LogPanicf(format string, v ...interface{}) {
	if logLevel <= PANIC {
		format = "[PANIC] " + format
		s := format + fmt.Sprintf(format, v...)
		logger.Output(2, s)
		panic(s)
	}
}

func LogPanic(format string, v ...interface{}) {
	if logLevel <= PANIC {
		s := "[PANIC] " + fmt.Sprintln(v...)
		logger.Output(2, s)
		panic(s)
	}
}

func LogFatalf(format string, v ...interface{}) {
	if logLevel <= FATAL {
		format = "[FATAL] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

func LogFatal(v ...interface{}) {
	if logLevel <= FATAL {
		s := "[FATAL] " + fmt.Sprintln(v...)
		logger.Output(2, s)
		os.Exit(1)
	}
}

func LogErrorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		format = "[ERROR] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogError(v ...interface{}) {
	if logLevel <= ERROR {
		s := "[ERROR] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func LogWarnf(format string, v ...interface{}) {
	if logLevel <= WARNING {
		format = "[WARNING] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogWarn(v ...interface{}) {
	if logLevel <= WARNING {
		s := "[WARNING] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func LogInfof(format string, v ...interface{}) {
	if logLevel <= INFO {
		format = "[INFO] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogInfo(v ...interface{}) {
	if logLevel <= INFO {
		s := "[INFO] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}

func LogDebugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		format = "[DEBUG] " + format
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func LogDebug(v ...interface{}) {
	if logLevel <= DEBUG {
		s := "[DEBUG] " + fmt.Sprintln(v...)
		logger.Output(2, s)
	}
}
