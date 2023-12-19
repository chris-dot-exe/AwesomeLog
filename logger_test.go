package log

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type foo struct {
	Foo    string
	Bar    string
	Foobar struct {
		Meeps []string
	}
}

func logs(lvl LogLevel) string {

	return Sprintln(lvl, lvl.String(), " Println")
}

func TestLogLevel(t *testing.T) {
	fmt.Println("For manual level test:")

	cfg := DefaultLevelConfig()
	cfg.Debug = LevelConfig{
		ShowLineNumber:   false,
		ShowFunctionName: true,
		ShowFilePath:     false,
	}

	SetLevelConfig(cfg)

	for name, lvl := range level {
		SetLogLevel(lvl)
		ShowColors(true)
		ShowColorsInLogs(true)

		fmt.Println("Level set to ", name)
		msg := "Test Message"
		Println(VERBOSE, msg)
		Println(DEBUG, msg)
		Println(INFO, msg)
		Println(WARN, msg)
		Println(ERROR, msg)
		Println(CRITICAL, msg)
		Println(NONE, msg)
		fmt.Println("----------------------")
	}
	fmt.Println("END manual level test")
	fmt.Println("----------------------")
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel(DEBUG)
	ShowTimestamp(false)
	ShowColors(true)
	ShowColorsInLogs(true)
	SetLevelConfig(DefaultLevelConfig())

	expexted := ANSI_RED_BACKGROUND + ANSI_WHITE + "[DEBUG]" + ANSI_RESET +
		"[logger_test.go:logs:21] DEBUG Println\n" +
		ANSI_YELLOW_BACKGROUND + ANSI_BLACK + "[WARN]" + ANSI_RESET +
		" WARN Println\n" +
		ANSI_BLUE_BACKGROUND + ANSI_WHITE + "[INFO]" + ANSI_RESET +
		" INFO Println\n"

	output := ""
	output += logs(VERBOSE)
	output += logs(DEBUG)
	output += logs(WARN)
	output += logs(INFO)
	output += logs(NONE)
	//fmt.Print(output)

	if output != expexted {
		t.Error("output is not as expected")
	}

}

func TestSetLogLevelByString(t *testing.T) {
	SetLogLevelByString("DEBUG")
	ShowTimestamp(false)
	ShowColorsInLogs(true)
	SetLevelConfig(DefaultLevelConfig())

	expected := ANSI_RED_BACKGROUND + ANSI_WHITE + "[DEBUG]" + ANSI_RESET +
		"[logger_test.go:logs:21] DEBUG Println\n" +
		ANSI_YELLOW_BACKGROUND + ANSI_BLACK + "[WARN]" + ANSI_RESET +
		" WARN Println\n" +
		ANSI_BLUE_BACKGROUND + ANSI_WHITE + "[INFO]" + ANSI_RESET +
		" INFO Println\n"

	output := ""
	output += logs(VERBOSE)
	output += logs(DEBUG)
	output += logs(WARN)
	output += logs(INFO)
	output += logs(NONE)
	//fmt.Print(output)

	if output != expected {
		t.Error("output is not as expected")
		fmt.Println("Expected: ")
		fmt.Print(expected)
		fmt.Println("Got: ")
		fmt.Print(output)
	}
}

func TestPrettyPrint(t *testing.T) {
	ShowColorsInLogs(true)
	SetLogLevel(INFO)
	demo := foo{
		Foo: "Test",
		Bar: "Test",
		Foobar: struct {
			Meeps []string
		}{[]string{"Meep", "Meep2", "Meep2.1"}},
	}

	expected := ANSI_BLUE_BACKGROUND + ANSI_WHITE + "[INFO]" + ANSI_RESET +
		" [\n" +
		"  {\n" +
		"    \"Foo\": \"Test\",\n" +
		"    \"Bar\": \"Test\",\n" +
		"    \"Foobar\": {\n" +
		"      \"Meeps\": [\n" +
		"        \"Meep\",\n" +
		"        \"Meep2\",\n" +
		"        \"Meep2.1\"\n" +
		"      ]\n" +
		"    }\n" +
		"  }\n" +
		"]\n"

	output := SprettyPrint(INFO, demo)

	if output != expected {
		t.Error("output is not as expected")

		fmt.Println("Expected: ")
		fmt.Print(expected)
		fmt.Println("Got: ")
		fmt.Print(output)
	}
}

func TestShowColorsInLogsActive(t *testing.T) {
	SetLogLevel(DEBUG)
	ShowColors(true)
	ShowColorsInLogs(true)
	ShowTimestamp(false)
	ShowCaller(false)
	expected := ANSI_RED_BACKGROUND + ANSI_WHITE + "[DEBUG]" + ANSI_RESET + " test\n"
	output := Sprintln(DEBUG, "test")

	if expected != output {
		t.Error("output is not as expected")

		fmt.Println("Expected: ")
		fmt.Print(expected)
		fmt.Println("Got: ")
		fmt.Print(output)
	}
}

func TestShowColorsInLogsInactive(t *testing.T) {
	SetLogLevel(DEBUG)
	ShowColors(false)
	ShowTimestamp(false)
	ShowCaller(false)
	expected := "[DEBUG] test\n"
	output := Sprintln(DEBUG, "test")

	if expected != output {
		t.Error("output is not as expected")

		fmt.Println("Expected: ")
		fmt.Print(expected)
		fmt.Println("Got: ")
		fmt.Print(output)
	}
}

func TestSprintf(t *testing.T) {
	SetLogLevel(DEBUG)
	ShowColors(false)
	ShowTimestamp(false)
	ShowCaller(false)

	expected := "[DEBUG] test 42, 0.1\n"
	output := Sprintf(DEBUG, "test %d, %.1f\n", 42, 0.1)

	Printf(DEBUG, "test %d, %.1f\n", 42, 0.1)

	if expected != output {
		t.Error("output is not as expected")

		fmt.Println("Expected: ")
		fmt.Print(expected)
		fmt.Println("Got: ")
		fmt.Print(output)
	}
}

func TestLogWithCustomTimeFormat(t *testing.T) {
	SetLogLevel(DEBUG)
	ShowCaller(false)
	ShowTimestamp(true)

	format := "2006/01/02 15:04:05.000000"
	now := time.Now()
	SetTimeFormat(format)

	expected := fmt.Sprintf("%s [DEBUG] Hello, World!", now.Format(format))

	msg := Message{
		Time:    now,
		Level:   DEBUG,
		Caller:  Caller{},
		Message: "Hello, World!",
	}

	result := stringify(msg)

	assert.Equal(t, expected, result)
}
