package log

import (
	log2 "log"
	"testing"
)

// customDiscard is a custom io.Writer that discards output.
// We use this instead of io.Discard because the Go standard library 'log' package
// has a hardcoded, highly optimized fast-path for io.Discard that skips all formatting and allocations.
// To make the benchmark between AwesomeLog and stdlib fair, we force the stdlib to format the string.
type customDiscard struct{}

func (customDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestBug_EmptyParametersPanic(t *testing.T) {
	SetOutput(customDiscard{})
	// Should not panic on empty params
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Panicked on empty parameters: %v", r)
		}
	}()
	Println()
	Print()
	Info()
}

func BenchmarkAwesomeLog_Info(b *testing.B) {
	SetOutput(customDiscard{})
	SetLogLevel(INFO)
	ShowTimestamp(false)
	ShowColorsInLogs(false)
	ShowCaller(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkAwesomeLog_Println(b *testing.B) {
	SetOutput(customDiscard{})
	SetLogLevel(INFO)
	ShowTimestamp(false)
	ShowColorsInLogs(false)
	ShowCaller(false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Println(INFO, "benchmark message")
	}
}

func BenchmarkStdlibLog_Println(b *testing.B) {
	log2.SetOutput(customDiscard{})
	log2.SetFlags(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log2.Println("benchmark message")
	}
}

func BenchmarkAwesomeLog_Info_Full(b *testing.B) {
	SetOutput(customDiscard{})
	SetLogLevel(INFO)
	ShowTimestamp(true)
	ShowColorsInLogs(false)
	ShowCaller(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}
}

func BenchmarkAwesomeLog_Println_Full(b *testing.B) {
	SetOutput(customDiscard{})
	SetLogLevel(INFO)
	ShowTimestamp(true)
	ShowColorsInLogs(false)
	ShowCaller(true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Println(INFO, "benchmark message")
	}
}

func BenchmarkStdlibLog_Println_Full(b *testing.B) {
	log2.SetOutput(customDiscard{})
	// LstdFlags = Ldate | Ltime. Wir fügen Lshortfile für die Caller Info hinzu.
	log2.SetFlags(log2.LstdFlags | log2.Lshortfile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log2.Println("benchmark message")
	}
}
