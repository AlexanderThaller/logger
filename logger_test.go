package logger

import (
	"bytes"
	"testing"
)

const (
	namet = name + ".Test"
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

func TestGetParentLevel(t *testing.T) {
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

func TestPrintMessage(t *testing.T) {
	l := New(namet + ".PrintMessage")

	m := [][]string{
		{"", "Test - Debug - "},
		{"Test", "Test - Debug - Test"},
		{"Test.Test", "Test - Debug - Test.Test"},
		{"Test.Test.Test", "Test - Debug - Test.Test.Test"},
	}

	r := getLogger("Test")
	r.Format = "{{.Logger}} - {{.Priority}} - {{.Message}}"
	r.NoColor = true

	for _, d := range m {
		l.Info("Checking: ", d)

		k := d[0]
		v := d[1]

		var b bytes.Buffer
		r.Output = &b

		printMessage(r, Debug, k)
		o := b.String()

		l.Debug("GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
		if o != v {
			l.Critical("GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
	}
}

func TestPrintColors(t *testing.T) {
	l := New("logger.Test.PrintColors")
	SetLevel("logger.Test.PrintColors", Disable)

	//TODO: Compare strings instead of printing.

	l.Debug("Debug")
	l.Info("Info")
	l.Notice("Notice")
	l.Warning("Warning")
	l.Error("Error")
	l.Critical("Critical")
	l.Alert("Alert")
	l.Emergency("Emergency")

	SetNoColor("logger.Test.PrintColors", true)
	l.Debug("NoColorDebug")
	l.Info("NoColorInfo")
	l.Notice("NoColorNotice")
	l.Warning("NoColorWarning")
	l.Error("NoColorError")
	l.Critical("NoColorCritical")
	l.Alert("NoColorAlert")
	l.Emergency("NoColorEmergency")
}

func TestCheckPriorityOK(t *testing.T) {
	l := New(namet + ".CheckPriority.OK")

	for k := range priorities {
		l.Info("Checking: ", k)

		e := checkPriority(k)
		l.Debug("Return of ", k, ": ", e)
		if e != nil {
			l.Critical(e)
			t.Fail()
		}
	}
}

func TestCheckPriorityFail(t *testing.T) {
	l := New(namet + ".CheckPriority.FAIL")

	k := Disable + 1

	l.Info("Checking: ", k)

	e := checkPriority(k)
	l.Debug("Return of ", k, ": ", e)
	if e == nil {
		l.Critical("Should not have succeeded")
		t.Fail()
		return
	}
}

func TestCheckPriorityFailDoesNotExist(t *testing.T) {
	l := New(namet + ".CheckPriority.FAIL.DoesNotExist")

	k := Disable + 1
	x := "priority does not exist"

	l.Info("Checking: ", k)

	e := checkPriority(k)
	l.Debug("Return of ", k, ": ", e)
	if e != nil {

		if e.Error() != x {
			l.Critical("Wrong error, EXPECTED: ", x, ", GOT: ", e.Error())
			t.Fail()
		}
	}
}

func BenchmarkLogRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log(".", Debug, "Test")
	}
}

func BenchmarkLogRootEmergency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log(".", Emergency, "Test")
	}
}

func BenchmarkLogChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log("BenchLogChild", Debug, "Test")
	}
}

func BenchmarkLogChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log("BenchLogChildChild.Test", Debug, "Test")
	}
}

func BenchmarkLogChildAllocated(b *testing.B) {
	SetLevel("BenchLogChildAllocated", Emergency)
	for i := 0; i < b.N; i++ {
		Log("BenchLogChildAllocated", Debug, "Test")
	}
}

func BenchmarkLogChildChildAllocated(b *testing.B) {
	SetLevel("BenchLogChildChildAllocated.Test", Emergency)
	for i := 0; i < b.N; i++ {
		Log("BenchLogChildChildAllocated.Test", Debug, "Test")
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
	l.Output = &a

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		printMessage(l, Debug, "Message")
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
