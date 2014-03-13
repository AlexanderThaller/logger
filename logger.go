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

// Logger represent different lognames and priorities. They can have
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

// Priority defines how important a log message is. Loggers will output
// messages which are above their priority level.
type Priority int

// Different priority levels ordered by their severity.
const (
	Debug Priority = iota
	Info
	Notice
	Warning
	Error
	Critical
	Alert
	Emergency
	Disable
)

// DefaultPriority of the root logger.
const (
	DefaultPriority Priority = Notice
)

const (
	defroot      = Logger(".")
	defseperator = "."
)

var (
	format     = "[{{.Time}} {{.Priority}} {{.Logger}}] - {{.Message}}.\n"
	timeformat = time.RFC3339

	priorities     map[Priority]string
	loggers        map[Logger]logger
	formattemplate template.Template
)

func init() {
	loggers = make(map[Logger]logger)
	l := logger{Logger: defroot,
		Priority:   DefaultPriority,
		Format:     Format(format),
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
	priorities[Disable] = "Disabled"
}

// New will return a logger with the given name.
func New(na string) (log Logger) {
	return Logger(na)
}

// SetLevel sets the priority level for the logger which should be
// logged.
func SetLevel(lo Logger, pr Priority) (err error) {
	err = errors.New("this priority does not exist")

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
		Format:     Format(format),
		TimeFormat: timeformat,
	}

	loggers[lo] = l
	return
}

// GetLevel returns the priority level of the given logger.
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

// ParsePriority tries to parse the priority by the given string.
func ParsePriority(pr string) (Priority, error) {
	for k, v := range priorities {
		if v == pr {
			return k, nil
		}
	}

	e := errors.New("Can not parse " + pr + ". Using DefaultPriority")
	return DefaultPriority, e
}

// NamePriority returns the string value of the given priority.
func NamePriority(pr Priority) string {
	return priorities[pr]
}

// SetFormat changes the message format for the given logger. Avaivable
// fields are:
//
// Time: The time when the message is printed.
//
// Logger: The name of the logger.
//
// Priority: The priority of the logger.
//
// Message: The output message.
//
// The default Format is:
//
// "[{{.Time}} {{.Logger}} {{.Priority}}] - {{.Message}}.\n"
func SetFormat(lo Logger, fo Format) (err error) {
	l := getLogger(lo)
	l.Format = fo
	loggers[lo] = l

	return
}

// SetTimeFormat sets the TimeFormat which will be used in the message
// format for the specified logger
//
// The default format is: RFC3339
func SetTimeFormat(lo Logger, fo string) (err error) {
	l := getLogger(lo)
	l.TimeFormat = fo
	loggers[lo] = l

	return
}

// Log logs a message using the given logger at a given priority.
func Log(lo Logger, pr Priority, me ...interface{}) {
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

// DebugM logs a message with the Debug priority.
func (lo Logger) Debug(me ...interface{}) {
	Log(lo, Debug, me...)
}

// InfoM logs a message with the Debug priority.
func (lo Logger) Info(me ...interface{}) {
	Log(lo, Info, me...)
}

// NoticeM logs a message with the Debug priority.
func (lo Logger) Notice(me ...interface{}) {
	Log(lo, Notice, me...)
}

// WarningM logs a message with the Debug priority.
func (lo Logger) Warning(me ...interface{}) {
	Log(lo, Warning, me...)
}

// ErrorM logs a message with the Debug priority.
func (lo Logger) Error(me ...interface{}) {
	Log(lo, Error, me...)
}

// CriticalM logs a message with the Debug priority.
func (lo Logger) Critical(me ...interface{}) {
	Log(lo, Critical, me...)
}

// AlertM logs a message with the Debug priority.
func (lo Logger) Alert(me ...interface{}) {
	Log(lo, Alert, me...)
}

// EmergencyM logs a message with the Debug priority.
func (lo Logger) Emergency(me ...interface{}) {
	Log(lo, Emergency, me...)
}
