package logger

import (
	"errors"
	"log"
)

const (
	LvlError = 0
	LvlInfo  = 1
	LvlDebug = 2
)

var lvl = 1

func SetLevel(lv int) (er error) {
	switch lv {
	case LvlError:
		lvl = LvlError
	case LvlInfo:
		lvl = LvlInfo
	case LvlDebug:
		lvl = LvlDebug
	default:
		er = errors.New("Not a valid level.")
		return
	}

	return
}

// Print a message in the style ERROR: Message if the lvl is or is bigger than
// LvlError
func Error(ms string) (ou string) {
	// Return if lvl is below error
	if lvl < LvlError {
		return
	}

	p := "ERROR: "

	if len(ms) == 0 {
		p = "ERROR"
	}

	ou = p + ms
	log.Print(ou)

	return
}

// Print a message in the style INFO : Message if the lvl is or is bigger than
// LvlInfo
func Info(ms string) (ou string) {
	// Return if lvl is below info
	if lvl < LvlInfo {
		return
	}

	p := "INFO : "

	if len(ms) == 0 {
		p = "INFO"
	}

	ou = p + ms
	log.Print(ou)

	return
}

// Print a message in the style DEBUG: Message if the lvl is or is bigger than
// LvlDebug
func Debug(ms string) (ou string) {
	// Return if lvl is below debug
	if lvl < LvlDebug {
		return
	}

	p := "DEBUG: "

	if len(ms) == 0 {
		p = "DEBUG"
	}

	ou = p + ms
	log.Print(ou)

	return
}
