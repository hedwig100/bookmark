package slog

import (
	"fmt"
	"log"
	"os"
)

var (
	logErr   *log.Logger
	logWarn  *log.Logger
	logInfo  *log.Logger
	logDebug *log.Logger
)

func init() {
	logErr = log.New(os.Stderr, "\x1b[31m[ERROR]\x1b[0m", log.LstdFlags)
	logWarn = log.New(os.Stderr, "\x1b[33m[WARN]\x1b[0m", log.LstdFlags)
	logInfo = log.New(os.Stderr, "[INFO]", log.LstdFlags)
	logDebug = log.New(os.Stderr, "[DEBUG]", log.LstdFlags)
}

func Fatal(v ...any) {
	logErr.Println(fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...any) {
	logErr.Println(fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Err(v ...any) {
	logErr.Println(fmt.Sprint(v...))
}

func Errf(format string, v ...any) {
	logErr.Println(fmt.Sprintf(format, v...))
}

func Warn(v ...any) {
	logWarn.Println(fmt.Sprint(v...))
}

func Warnf(format string, v ...any) {
	logWarn.Println(fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	logInfo.Println(fmt.Sprint(v...))
}

func Infof(format string, v ...any) {
	logInfo.Println(fmt.Sprintf(format, v...))
}

func Debug(v ...any) {
	logDebug.Println(fmt.Sprint(v...))
}

func Debugf(format string, v ...any) {
	logDebug.Println(fmt.Sprintf(format, v...))
}
