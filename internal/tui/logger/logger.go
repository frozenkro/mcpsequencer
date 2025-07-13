package logger

import "log"

var Logger *log.Logger

func SetLogger(logger *log.Logger) {
	Logger = logger
}
