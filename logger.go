// Package logger provides a logging framework similar to those of python
// and haskell.
package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

// Format represents the format which will be used to print the message
// for an logger.
type Format string

// Loggers represent different lognames and priorities. They can have
// parents and child loggers which will inherit the priority of the
// parent if it has none. The hirachy of loggers is represented through
// sepperation with dots ('.'). The root logger has the name '.'.
type Logger string

type message struct {
	Time string
	Logger
	Priority string
	Message  string
}

type logger struct {
	Logger
	Priority
	Format
	TimeFormat string
}

// Priorities define how important a log message is. Loggers will output
// messages which are above their priority level.
type Priority int

const (
	Debug     Priority = iota // Debugging messages
	Info                      // Information
	Notice                    // Normal messages
	Warning                   // General Warning
	Error                     // General Errors
	Critical                  // Severe situations
	Alert                     // Take immediate action
	Emergency                 // System is unusable
)

const (
	defroot         = Logger(".")
	defseperator    = "."
	DefaultPriority = Notice
)

var (
	format     Format = "[{{.Time}} {{.Logger}} {{.Priority}}] - {{.Message}}.\n"
	timeformat        = time.RFC3339

	priorities     map[Priority]string
	loggers        map[Logger]logger
	formattemplate template.Template
)

func init() {
	loggers = make(map[Logger]logger)
	l := logger{Logger: defroot,
		Priority:   DefaultPriority,
		Format:     format,
		TimeFormat: timeformat,
	}

	loggers[defroot] = l

	priorities = make(map[Priority]string)
	priorities[Debug] = "Debug"
	priorities[Info] = "Info"
	priorities[Notice] = "Notice"
	priorities[Warning] = "Warning"
	priorities[Error] = "Error"
	priorities[Critical] = "Critical"
	priorities[Alert] = "Alert"
	priorities[Emergency] = "Emergency"
}

// Set the priority level of the logger.
func SetLevel(lo Logger, pr Priority) (err error) {
	err = errors.New("This priority does not exist")

	for k := range priorities {
		if k == pr {
			err = nil
			break
		}
	}

	if err != nil {
		return
	}

	l := logger{Logger: lo,
		Priority:   pr,
		Format:     format,
		TimeFormat: timeformat,
	}

	loggers[lo] = l
	return
}

// Get the priority level of the logger.
func GetLevel(lo Logger) (pri Priority) {
	l, e := loggers[lo]
	if e {
		pri = l.Priority
	} else {
		pri = getParentLevel(lo)
	}

	return
}

func getParentLevel(lo Logger) (pri Priority) {
	p := getParent(lo)
	pri = GetLevel(p)

	return
}

func getParent(lo Logger) (log Logger) {
	// Return root if root
	if lo == defroot {
		log = defroot
		return
	}

	s := strings.Split(string(lo), defseperator)

	// Return root if first level logger
	if len(s) == 1 {
		log = defroot
		return
	}

	// Return root if parent is empty
	if s[0] == "" {
		log = defroot
		return
	}

	l := len(s) - 1
	z := s[0:l]

	log = Logger(strings.Join(z, defseperator))

	return
}

func getLogger(lo Logger) (log logger) {
	l, e := loggers[lo]
	if e {
		log = l
	} else {
		log = getParentLogger(lo)
	}

	return
}

func getParentLogger(lo Logger) (log logger) {
	l := getParent(lo)
	log = getLogger(l)
	log.Logger = lo

	return
}

// Parses and returns the given priority.
func ParsePriority(pr string) Priority {
	for k, v := range priorities {
		if v == pr {
			return k
		}
	}

	e := errors.New("Can not parse " + pr + ". Using DefaultPriority")
	ErrorM("logger.ParsePriority", e)
	return DefaultPriority
}

// Returns the name of an priority.
func NamePriority(pr Priority) string {
	return priorities[pr]
}

// Set output format. Avaivable fields are:
// Time: The time when the message is printed.
// Logger: The name of the logger.
// Priority: The priority of the logger.
// Message: The output message.
// The default Format is:
// "[{{.Time}} {{.Logger}} {{.Priority}}] - {{.Message}}.\n"
func SetFormat(lo Logger, fo Format) (err error) {
	l := getLogger(lo)
	l.Format = fo
	loggers[lo] = l

	return
}

// Set the TimeFormat which will be used in the message format for the
// specified logger
// The default format is: RFC3339
func SetTimeFormat(lo Logger, fo string) (err error) {
	l := getLogger(lo)
	l.TimeFormat = fo
	loggers[lo] = l

	return
}

// Log a message using the given logger at a given priority.
func LogM(lo Logger, pr Priority, me ...interface{}) {
	l := getLogger(lo)

	if l.Priority > pr {
		return
	}

	printMessage(l, pr, os.Stderr, me...)
}

func printMessage(lo logger, pr Priority, wr io.Writer, me ...interface{}) {
	m := new(message)
	m.Time = time.Now().Format(string(lo.TimeFormat))
	m.Logger = lo.Logger
	m.Priority = priorities[pr]
	m.Message = fmt.Sprint(me...)

	s := formatMessage(m, lo.Format)

	fmt.Fprint(wr, s)
}

func formatMessage(me *message, fo Format) (so string) {
	so = strings.Replace(string(fo), "{{.Time}}", me.Time, -1)
	so = strings.Replace(so, "{{.Logger}}", string(me.Logger), -1)
	so = strings.Replace(so, "{{.Priority}}", me.Priority, -1)
	so = strings.Replace(so, "{{.Message}}", me.Message, -1)

	return
}

// Log a message at Debug priority.
func DebugM(lo Logger, me ...interface{}) {
	LogM(lo, Debug, me...)
}

// Log a message at Info priority.
func InfoM(lo Logger, me ...interface{}) {
	LogM(lo, Info, me...)
}

// Log a message at Notice priority.
func NoticeM(lo Logger, me ...interface{}) {
	LogM(lo, Notice, me...)
}

// Log a message at Warning priority.
func WarningM(lo Logger, me ...interface{}) {
	LogM(lo, Warning, me...)
}

// Log a message at Error priority.
func ErrorM(lo Logger, me ...interface{}) {
	LogM(lo, Error, me...)
}

// Log a message at Critical priority.
func CriticalM(lo Logger, me ...interface{}) {
	LogM(lo, Critical, me...)
}

// Log a message at Alert priority.
func AlertM(lo Logger, me ...interface{}) {
	LogM(lo, Alert, me...)
}

// Log a message at Emergency priority.
func EmergencyM(lo Logger, me ...interface{}) {
	LogM(lo, Emergency, me...)
}
