package log

import (
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
			err:   log.New(os.Stderr, "[ERR]: ", logFlags),
			fatal: log.New(os.Stderr, "[FATAL]: ", logFlags),
		}
	})
	Instance.level = constants.Err
	Instance.initialised = true
}

func (l *Logger) SetLogLevel(level int) {
	l.level = level
}

func (l *Logger) Debug(message string) {
	if l.level > constants.Debug {
		return
	}

	message = filterMessageForLog(message)
	l.debug.Println(message)
}

func (l *Logger) Info(message string) {
	if l.level > constants.Info {
		return
	}

	message = filterMessageForLog(message)
	l.info.Println(message)
}

func (l *Logger) Warn(message string) {
	if l.level > constants.Warn {
		return
	}

	message = filterMessageForLog(message)
	l.warn.Println(message)
}

func (l *Logger) Err(message string, err error) {
	if l.level > constants.Err {
		return
	}

	message = filterMessageForLog(message)
	if err != nil {
		l.err.Printf("%v, err: %v", message, err.Error())
		return
	}

	l.err.Println(message)
}

func (l *Logger) Fatal(message string, err error) {
	if l.level > constants.Fatal {
		return
	}

	message = filterMessageForLog(message)

	if err != nil {
		l.fatal.Printf("%v, err: %v", message, err.Error())
	} else {
		l.fatal.Println(message)
	}

	os.Exit(1)
}

func filterMessageForLog(message string) string {
	r := regexp.MustCompile(`/password[^,{}\]\[]*/gim`)
	message = r.ReplaceAllString(message, "password:'********'")
	return message
}
