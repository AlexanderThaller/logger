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
	TimeFormat     string
	FormatTemplate *template.Template
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
	defroot      = Logger(".")
	defseperator = "."
	defpriority  = Notice
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
		Priority:   defpriority,
		Format:     format,
		TimeFormat: timeformat,
	}
	t := template.New("FormatTemplate")
	t, _ = t.Parse(string(format))
	l.FormatTemplate = t

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

	t := template.New("FormatTemplate")
	t, _ = t.Parse(string(format))
	log.FormatTemplate = t

	return
}

// Set output format. Avaivable fields are:
// Time: The time when the message is printed.
// Logger: The name of the logger.
// Priority: The priority of the logger.
// Message: The output message.
// The default Format is:
// "[{{.Time}} {{.Logger}} {{.Priority}}] - {{.Message}}.\n"
func SetFormat(lo Logger, fo Format) (err error) {
	t := template.New("FormatTemplate")
	t, err = t.Parse(string(fo))
	if err != nil {
		return
	}

	l := getLogger(lo)
	l.Format = fo
	l.FormatTemplate = t
	loggers[lo] = l

	return
}

// Set the timeformat that will be used to format the time in the format.
// The default format is: RFC3339
func SetTimeFormat(fo string) (err error) {
	timeformat = fo
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

	e := lo.FormatTemplate.Execute(wr, m)
	if e != nil {
		fmt.Println(e)
	}
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
