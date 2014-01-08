package logger

import (
	"bytes"
	"testing"
)

func TestGetLevel(t *testing.T) {
	n := New("logger.Test.GetLevel")

	n.Info(n, "Starting")
	m := make(map[Logger]Priority)
	m[""] = DefaultPriority
	m["."] = DefaultPriority
	m["Test"] = DefaultPriority
	m[".Test"] = DefaultPriority

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
			n.Error(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
		n.Debug(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	n.Info(n, "Finished")
}

func TestgetParentLevel(t *testing.T) {
	n := New("logger.Test.getParentLevel")

	n.Info(n, "Starting")
	m := make(map[Logger]Priority)
	m["."] = DefaultPriority
	m["Test"] = DefaultPriority
	m["Test.Test"] = DefaultPriority

	SetLevel("Test2", Emergency)
	m["Test2"] = DefaultPriority
	m["Test2.Test"] = Emergency

	for k, v := range m {
		o := getParentLevel(k)
		if o != v {
			n.Error(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
		n.Debug(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	n.Info(n, "Finished")
}

func TestgetParent(t *testing.T) {
	n := New("logger.Test.getParent")

	n.Info(n, "Starting")
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
			n.Error(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
		n.Debug(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	n.Info(n, "Finished")
}

func TestprintMessage(t *testing.T) {
	n := New("logger.Test.printMessage")

	n.Info(n, "Starting")
	m := [][]string{
		{"", "Test - Debug - "},
		{"Test", "Test - Debug - Test"},
		{"Test.Test", "Test - Debug - Test.Test"},
		{"Test.Test.Test", "Test - Debug - Test.Test.Test"},
	}

	SetFormat("Test", "{{.Logger}} - {{.Priority}} - {{.Message}}")
	l := getLogger("Test")

	for i := range m {
		a := m[i]

		k := a[0]
		v := a[1]

		var b bytes.Buffer
		printMessage(l, Debug, &b, k)
		o := b.String()
		if o != v {
			n.Error(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
		n.Debug(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}

	n.Info(n, "Finished")
}

func BenchmarkLogMRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log(".", Debug, "Test")
	}
}

func BenchmarkLogMRootEmergency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log(".", Emergency, "Test")
	}
}

func BenchmarkLogMChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log("BenchLogMChild", Debug, "Test")
	}
}

func BenchmarkLogMChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log("BenchLogMChildChild.Test", Debug, "Test")
	}
}

func BenchmarkLogMChildAllocated(b *testing.B) {
	SetLevel("BenchLogMChildAllocated", Emergency)
	for i := 0; i < b.N; i++ {
		Log("BenchLogMChildAllocated", Debug, "Test")
	}
}

func BenchmarkLogMChildChildAllocated(b *testing.B) {
	SetLevel("BenchLogMChildChildAllocated.Test", Emergency)
	for i := 0; i < b.N; i++ {
		Log("BenchLogMChildChildAllocated.Test", Debug, "Test")
	}
}

func BenchmarkGetParentRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent(".")
	}
}

func BenchmarkGetParentChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("BenchgetParentChild")
	}
}

func BenchmarkGetParentChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("BenchgetParentChildChild.Test")
	}
}

func BenchmarkGetParentChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("BenchgetParentChildChild.Test.Test")
	}
}

func BenchmarkGetParentChildChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("BenchgetParentChildChildChild.Test.Test")
	}
}

func BenchmarkGetParentChildChildChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("BenchgetParentChildChildChildChild.Test.Test.Test")
	}
}

func BenchmarkPrintMessage(b *testing.B) {
	var a bytes.Buffer
	l := getLogger("BenchprintMessage")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		printMessage(l, Debug, &a, "Message")
	}
}

func BenchmarkFormatMessage(b *testing.B) {
	l := getLogger("BenchformatMessage")

	m := new(message)
	m.Time = "Mo 30 Sep 2013 20:29:19 CEST"
	m.Logger = l.Logger
	m.Priority = "Debug"
	m.Message = "Test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatMessage(m, l.Format)
	}
}
