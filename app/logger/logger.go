package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	Info        *log.Logger
	Debug       *log.Logger
	Warn        *log.Logger
	Err         *log.Logger
	Fatal       *log.Logger
	level       severity
	initialised bool
}

type severity int

var defaultLogger *Logger

const logFlags = log.Ldate | log.Ltime

// Severity levels.
const (
	sDebug severity = iota
	sInfo
	sWarn
	sErr
	sFatal
)

func Init() {
	initialise()
}

func initialise() {
	defaultLogger = &Logger{
		Debug: log.New(os.Stderr, "DEBUG: ", logFlags),
		Info:  log.New(os.Stderr, "INFO : ", logFlags),
		Warn:  log.New(os.Stderr, "WARNING: ", logFlags),
		Err:   log.New(os.Stderr, "ERROR: ", logFlags),
		Fatal: log.New(os.Stderr, "FATAL: ", logFlags),
	}
	defaultLogger.level = sErr
	defaultLogger.initialised = true
}

func (logger *Logger) output(level severity, message string) {
	switch level {
	case sDebug:
		logger.Debug.Output(3, message)
	case sInfo:
		logger.Info.Output(3, message)
	case sWarn:
		logger.Warn.Output(3, message)
	case sErr:
		logger.Err.Output(3, message)
	case sFatal:
		logger.Fatal.Output(3, message)
	default:
		panic(fmt.Sprintln("Invlid log severity: ", level))
	}
}

func SetLogLevel(level severity) {
	defaultLogger.level = level
	fmt.Printf("Log level set to %d\n", level)
}

func DEBUG(message string) {
	if defaultLogger.level <= sDebug {
		defaultLogger.output(sDebug, message)
	}
}

func INFO(message string) {
	if defaultLogger.level <= sInfo {
		defaultLogger.output(sInfo, message)
	}
}

func WARN(message string) {
	if defaultLogger.level <= sWarn {
		defaultLogger.output(sWarn, message)
	}
}

func ERR(message string) {
	if defaultLogger.level <= sErr {
		defaultLogger.output(sErr, message)
	}
}

func FATAL(message string) {
	if defaultLogger.level <= sFatal {
		defaultLogger.output(sFatal, message)
	}
}

// func messageBuilder(msg string) string {
// 	r, _ := regexp.Compile("/password[^,}\]]*/gim")
// 	msg = r.ReplaceAllString(msg, "password:'********'")
// 	return msg
// }
