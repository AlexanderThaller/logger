package logger

import (
	"bytes"
	"testing"
)

func Test_GetLevel(t *testing.T) {
	n := Logger("logger.Test.GetLevel")

	InfoM(n, "Starting")
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
		DebugM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	InfoM(n, "Finished")
}

func Test_getParentLevel(t *testing.T) {
	n := Logger("logger.Test.getParentLevel")

	InfoM(n, "Starting")
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
		DebugM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	InfoM(n, "Finished")
}

func Test_getParent(t *testing.T) {
	n := Logger("logger.Test.getParent")

	InfoM(n, "Starting")
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
		DebugM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}
	InfoM(n, "Finished")
}

func Test_printMessage(t *testing.T) {
	n := Logger("logger.Test.printMessage")

	InfoM(n, "Starting")
	m := [][]string{
		{"", "Test - Debug - "},
		{"Test", "Test - Debug - Test"},
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
			ErrorM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
			t.Fail()
		}
		DebugM(n, "GOT: '", o, "', EXPECED: '", v, "'", ", KEY: '", k, "'")
	}

	InfoM(n, "Finished")
}

func Benchmark_LogMRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LogM(".", Debug, "Test")
	}
}

func Benchmark_LogMRootEmergency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LogM(".", Emergency, "Test")
	}
}

func Benchmark_LogMChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LogM("Bench_LogMChild", Debug, "Test")
	}
}

func Benchmark_LogMChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LogM("Bench_LogMChildChild.Test", Debug, "Test")
	}
}

func Benchmark_LogMChildAllocated(b *testing.B) {
	SetLevel("Bench_LogMChildAllocated", Emergency)
	for i := 0; i < b.N; i++ {
		LogM("Bench_LogMChildAllocated", Debug, "Test")
	}
}

func Benchmark_LogMChildChildAllocated(b *testing.B) {
	SetLevel("Bench_LogMChildChildAllocated.Test", Emergency)
	for i := 0; i < b.N; i++ {
		LogM("Bench_LogMChildChildAllocated.Test", Debug, "Test")
	}
}

func Benchmark_getParentRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent(".")
	}
}

func Benchmark_getParentChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("Bench_getParentChild")
	}
}

func Benchmark_getParentChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("Bench_getParentChildChild.Test")
	}
}

func Benchmark_getParentChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("Bench_getParentChildChild.Test.Test")
	}
}

func Benchmark_getParentChildChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("Bench_getParentChildChildChild.Test.Test")
	}
}

func Benchmark_getParentChildChildChildChildChild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getParent("Bench_getParentChildChildChildChild.Test.Test.Test")
	}
}

func Benchmark_printMessage(b *testing.B) {
		var a bytes.Buffer
	  l := getLogger("Bench_printMessage")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
      printMessage(l, Debug, &a, "Message")
    }
}

func Benchmark_printMessageExecute(b *testing.B) {
  var a bytes.Buffer
  l := getLogger("Bench_printMessageExecute")

	m := new(message)
  m.Time = "Mo 30 Sep 2013 20:29:19 CEST"
	m.Logger = l.Logger
	m.Priority = "Debug"
	m.Message = "Test"

  b.ResetTimer()
  for i := 0; i < b.N; i++ {
    l.FormatTemplate.Execute(&a, m)
  }
}
