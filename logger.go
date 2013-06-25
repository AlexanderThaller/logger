package logger

import (
	"errors"
	"fmt"
	"log"
	"misc"
)

const (
	LvlError = 0
	LvlInfo  = 1
	LvlDebug = 2
)

var lvl = 1

// Set the level of the logger.
func SetLevel(lv int) (err error) {
	switch lv {
	case LvlError:
		lvl = LvlError
	case LvlInfo:
		lvl = LvlInfo
	case LvlDebug:
		lvl = LvlDebug
	default:
		err = errors.New("Not a valid level.")
		return
	}

	return
}

// Sets the output to use the specified logfile.
func SetLogfile(pa string) (err error) {
	f, err := misc.AppendFile(pa)
	if err != nil {
		return
	}

	log.SetOutput(f)

	return
}

// Print a message in the style ERROR: Message if the lvl is or is bigger than
// LvlError
func Error(v ...interface{}) (out string) {
	// Return if lvl is below error
	if lvl < LvlError {
		return
	}

	p := "ERROR: "
	ms := fmt.Sprint(v...)

	if len(ms) == 0 {
		p = "ERROR"
	}

	out = p + ms
	log.Print(out)

	return
}

// Print a message in the style INFO : Message if the lvl is or is bigger than
// LvlInfo
func Info(v ...interface{}) (out string) {
	// Return if lvl is below info
	if lvl < LvlInfo {
		return
	}

	p := "INFO : "
	ms := fmt.Sprint(v...)

	if len(ms) == 0 {
		p = "INFO"
	}

	out = p + ms
	log.Print(out)

	return
}

// Print a message in the style DEBUG: Message if the lvl is or is bigger than
// LvlDebug
func Debug(v ...interface{}) (out string) {
	// Return if lvl is below debug
	if lvl < LvlDebug {
		return
	}

	p := "DEBUG: "
	ms := fmt.Sprint(v...)

	if len(ms) == 0 {
		p = "DEBUG"
	}

	out = p + ms
	log.Print(out)

	return
}
