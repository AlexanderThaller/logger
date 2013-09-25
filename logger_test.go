package logger

import (
	"testing"
)

func init() {
	SetLevel("logger.Test", Debug)
}

func Test_GetLevel(t *testing.T) {
	n := Logger("logger.Test.GetLevel")

	InfoM(n, "Starting test")
	m := make(map[Logger]Priority)
	m[""] = defpriority
	m["."] = defpriority
	m["Test"] = defpriority
	m[".Test"] = defpriority

	SetLevel("Test2", Emergency)
	m["Test2"] = Emergency
	m["Test2.Test"] = Emergency
	m["Test2.Test.Test"] = Emergency
	m["Test2.Test.Test.Test"] = Emergency
	m["Test2.Test.Test.Test.Test"] = Emergency
	m["Test2.Test.Test.Test.Test.Test"] = Emergency

	for k, v := range m {
		o := GetLevel(k)
		if o != v {
			ErrorM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
	}
}

func Test_getParentLevel(t *testing.T) {
	n := Logger("logger.Test.getParentLevel")

	InfoM(n, "Starting test")
	m := make(map[Logger]Priority)
	m["."] = defpriority
	m["Test"] = defpriority
	m["Test.Test"] = defpriority

	SetLevel("Test2", Emergency)
	m["Test2"] = defpriority
	m["Test2.Test"] = Emergency

	for k, v := range m {
		o := getParentLevel(k)
		if o != v {
			ErrorM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
	}
}

func Test_getParent(t *testing.T) {
	n := Logger("logger.Test.getParent")

	InfoM(n, "Starting test")
	m := [][]Logger{
		{"", "."},
		{".Test", "."},
		{".", "."},
		{"Test", "."},
		{"Test.Test", "Test"},
		{"Test.Test.Test", "Test.Test"},
		{"Test.Test.Test.Test", "Test.Test.Test"},
	}

	for i := range m {
		a := m[i]

		k := a[0]
		v := a[1]

		o := getParent(k)
		if o != v {
			ErrorM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
	}
}
