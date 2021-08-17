package log

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/MDAkramSiddiqui/sf-covid-api/app/constants"
)

type Logger struct {
	debug       *log.Logger
	info        *log.Logger
	warn        *log.Logger
	err         *log.Logger
	fatal       *log.Logger
	level       int
	initialised bool
}

const logFlags = log.Ldate | log.Ltime

var Instance *Logger
var logSyncOnce sync.Once

func Init() {
	initialise()
}

func initialise() {
	logSyncOnce.Do(func() {
		Instance = &Logger{
			debug: log.New(os.Stderr, "[DEBUG]: ", logFlags),
			info:  log.New(os.Stderr, "[INFO]: ", logFlags),
			warn:  log.New(os.Stderr, "[WARN]: ", logFlags),
			err:   log.New(os.Stderr, "[ERROR]: ", logFlags),
			fatal: log.New(os.Stderr, "[FATAL]: ", logFlags),
		}
	})
	Instance.level = constants.Err
	Instance.initialised = true
}

func (l *Logger) SetLogLevel(level int) {
	l.level = level
}

// print info debug messages
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.level > constants.Debug {
		return
	}

	message = buildMessageForLog(message, args...)
	l.debug.Println(message)
}

// print info log messages
func (l *Logger) Info(message string, args ...interface{}) {
	if l.level > constants.Info {
		return
	}

	message = buildMessageForLog(message, args...)
	l.info.Println(message)
}

// print warn log messages
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.level > constants.Warn {
		return
	}

	message = buildMessageForLog(message, args...)
	l.warn.Println(message)
}

// print error log messages
func (l *Logger) Err(message string, args ...interface{}) {
	if l.level > constants.Err {
		return
	}

	message = buildMessageForLog(message, args...)
	l.err.Println(message)
}

// print fatal log messages and exit app with status 1
func (l *Logger) Fatal(message string, args ...interface{}) {
	if l.level > constants.Fatal {
		return
	}

	message = buildMessageForLog(message, args...)
	l.fatal.Println(message)
	os.Exit(1)
}

// build messages for logs,
// masks password keys and interpolates messages
func buildMessageForLog(message string, args ...interface{}) string {
	r := regexp.MustCompile(`/password[^,{}\]\[]*/gim`)
	message = r.ReplaceAllString(message, "password: ********")
	message = fmt.Sprintf(message, args...)
	return message
}
