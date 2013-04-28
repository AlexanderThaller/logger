package logger

import (
	"log"
	"strconv"
	"testing"
)

func Test_SetLevel(t *testing.T) {
	m := map[int]string{
		-1:       "Not a valid level.",
		-100:     "Not a valid level.",
		LvlError: "",
		LvlInfo:  "",
		LvlDebug: "",
		3:        "Not a valid level.",
		4:        "Not a valid level.",
	}

	for k, v := range m {
		lvl := k

		var o string
		e := SetLevel(lvl)

		if e == nil {
			o = ""
		} else {
			o = e.Error()
		}

		if o != v {
			log.Print("Level: ", lvl)
			log.Print("GOT: '", o, "', EXPECED: '", v, "'")
			t.Fail()
		}
	}
}

func Test_Error(t *testing.T) {
	m := [][]string{
		{"ERROR: Test", "0", "Test"},
		{"ERROR: Test", "1", "Test"},
		{"ERROR: Test", "2", "Test"},
		{"ERROR", "2", ""},
	}

	for i := range m {
		a := m[i]

		lvl, _ := strconv.Atoi(a[1])
		msg := a[2]

		SetLevel(lvl)

		v := a[0]
		o := Error(msg)
		if o != v {
			log.Print("Level: ", lvl)
			log.Print("GOT: '", o, "', EXPECED: '", v, "'")
			t.Fail()
		}
	}
}

func Test_Info(t *testing.T) {
	m := [][]string{
		{"", "0", "Test"},
		{"INFO : Test", "1", "Test"},
		{"INFO : Test", "2", "Test"},
		{"INFO", "2", ""},
	}

	for i := range m {
		a := m[i]

		lvl, _ := strconv.Atoi(a[1])
		msg := a[2]

		SetLevel(lvl)

		v := a[0]
		o := Info(msg)
		if o != v {
			log.Print("Level: ", lvl)
			log.Print("GOT: '", o, "', EXPECED: '", v, "'")
			t.Fail()
		}
	}
}

func Test_Debug(t *testing.T) {
	m := [][]string{
		{"", "0", "Test"},
		{"", "1", "Test"},
		{"DEBUG: Test", "2", "Test"},
		{"DEBUG", "2", ""},
	}

	for i := range m {
		a := m[i]

		lvl, _ := strconv.Atoi(a[1])
		msg := a[2]

		SetLevel(lvl)

		v := a[0]
		o := Debug(msg)
		if o != v {
			log.Print("Level: ", lvl)
			log.Print("GOT: '", o, "', EXPECED: '", v, "'")
			t.Fail()
		}
	}
}
