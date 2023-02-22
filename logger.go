package log

import (
	"encoding/json"
	"errors"
	"fmt"
	log2 "log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/structs"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	config = DefaultLevelConfig()
}

// SetLogLevel defines to which LogLevel log messages should be shown.
//
// Default is VERBOSE
func SetLogLevel(lvl LogLevel) {
	logLevel = lvl
}

// ShowCaller defines if the caller (function name, line number, file path) should be shown on a global level.
func ShowCaller(show bool) {
	if config == nil {
		config = DefaultLevelConfig()
	}

	s := structs.New(config)

	for _, value := range s.Fields() {

		sfn := value.Field("ShowFunctionName")
		sfp := value.Field("ShowFilePath")
		sln := value.Field("ShowLineNumber")
		sfn.Set(show)
		sfp.Set(show)
		sln.Set(show)
	}
}

// ShowColors Defines if colored level tags should be shown in the console log.
func ShowColors(show bool) {
	showColors = show
}

// SetLogLevelByString defines to which LogLevel log messages should be shown based on the given string e.g. SetLogLevelByString("WARN")
// This is useful if the LogLevel is defined in a config file.
func SetLogLevelByString(lvlStr string) {
	lvlStr = strings.ToUpper(lvlStr)
	val, ok := level[lvlStr]
	if !ok {
		log2.Fatalf("LogLevel '%s' is not supported!\n", lvlStr)
		return
	}
	logLevel = val
}

// SetDefaultLevel defines which LogLevel should be used if no LogLevel is provided.
// This is useful if AwesomeLog is used as a drop-in replacement for the build-in log package in existing projects.
//
// Default is INFO
func SetDefaultLevel(lvl LogLevel) {
	defaultLevel = lvl
}

// ShowColorsInLogs if set to true colored level tags are always active.
// By default, colored level tags are only active when the log is written to a terminal
func ShowColorsInLogs(show bool) {
	colorsInLogs = show
}

// ShowTimestamp defines if the log message should be prefixed with a timestamp
func ShowTimestamp(show bool) {
	showTimestamp = show
}

// SetCallerMaxDepth set the max depth of the callers file path
func SetCallerMaxDepth(depth int) {
	maxDepthOfCallerPath = depth
}

// DefaultLevelConfig return the default level config for AwesomeLog
func DefaultLevelConfig() *Config {
	cfg := &Config{
		Verbose: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
			Handlers:         []Handler{log},
		},
		Debug: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
			Handlers:         []Handler{log},
		},
		Warn: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
			Handlers:         []Handler{log},
		},
		Info: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
			Handlers:         []Handler{log},
		},
	}
	return cfg
}

// SetLevelConfig set the config for AwesomeLog
func SetLevelConfig(cfg *Config) {
	config = cfg
}

// SetTimeFormat set the timeformat for log messages
func SetTimeFormat(format string) {
	timeFormat = format
}

// Println logs a message at the defined LogLevel a newline is appended
func Println(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	println(level, params...)
}

// Print logs a message at the defined LogLevel
func Print(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	logHandler(level, params...)
}

// Printf logs a message at the defined LogLevel and formats the message according to a format specifier
func Printf(paramsOriginal ...interface{}) {
	level, format, params := getLogLevel(true, paramsOriginal...)
	logHandler(level, fmt.Sprintf(format, params...))
}

// PrettyPrint logs a message at the defined LogLevel formatted as JSON
// Works only with exported fields.
func PrettyPrint(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	b, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		Fatal("unsupported input. error: ", err)
	}

	println(level, string(b))
}

func Sprintln(params ...interface{}) string {
	level, _, params := getLogLevel(false, params...)
	return sprintln(level, params...)
}

func Sprint(params ...interface{}) string {
	level, _, params := getLogLevel(false, params...)
	return sprint(level, params)
}

func Sprintf(paramsOriginal ...interface{}) string {
	level, format, params := getLogLevel(true, paramsOriginal...)
	return sprint(level, fmt.Sprintf(format, params...))
}

func SprettyPrint(params ...interface{}) string {
	level, _, params := getLogLevel(false, params...)
	b, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		Fatal("unsupported input. error: ", err)
	}

	return sprintln(level, string(b))
}

// Fatal calls log.Fatal of the built-in log package.
// This function is provided only for drop-in compatibility
func Fatal(params ...interface{}) {
	log2.Fatal(params...)
}

// Fatalf calls log.Fatalf of the built-in log package.
// This function is provided only for drop-in compatibility
func Fatalf(format string, params ...interface{}) {
	log2.Fatalf(format, params...)
}

// Fatalln calls log.Fatalln of the built-in log package.
// This function is provided only for drop-in compatibility
func Fatalln(params ...interface{}) {
	log2.Fatalln(params...)
}

// Panic calls log.Panic of the built-in log package.
// This function is provided only for drop-in compatibility
func Panic(params ...interface{}) {
	log2.Panic(params...)
}

// Panicf calls log.Panicf of the built-in log package.
// This function is provided only for drop-in compatibility
func Panicf(format string, params ...interface{}) {
	log2.Panicf(format, params...)
}

// Panicln calls log.Panicln of the built-in log package.
// This function is provided only for drop-in compatibility
func Panicln(params ...interface{}) {
	log2.Panicln(params...)
}

// stringify builds the log message string with colors and caller
func stringify(message Message) string {

	if config == nil {
		config = DefaultLevelConfig()
	}

	lvlName := strings.Title(strings.ToLower(message.Level.String()))

	s := structs.New(config)
	lvlField := s.Field(lvlName)

	cfg := lvlField.Value().(LevelConfig)

	prefix := ""
	caller := ""

	if showTimestamp {
		prefix = fmt.Sprintf("%s ", message.Time.Format(timeFormat))
	}

	if showColors && (colorsInLogs || isTerminal()) {
		prefix += fmt.Sprintf(message.Level.Color()+"[%s]"+ANSI_RESET, message.Level.String())
	} else {
		prefix += fmt.Sprintf("[%s]", message.Level.String())
	}

	if cfg.ShowFilePath || cfg.ShowFunctionName || cfg.ShowLineNumber {
		caller += "["
		if cfg.ShowFilePath {

			caller += fmt.Sprintf("%s:", message.Caller.Path)
		}
		if cfg.ShowFunctionName {
			caller += fmt.Sprintf("%s", message.Caller.FunctionName)
		}
		if cfg.ShowLineNumber {
			caller += fmt.Sprintf(":%d", message.Caller.LineNumber)
		}

		caller += "]"
	}
	return fmt.Sprintf("%s%s %s", prefix, caller, message.Message)
}

// buildMessage builds the Message object used by all log handlers
func buildMessage(level LogLevel, params ...interface{}) Message {
	now := time.Now()
	caller := Caller{}

	fpcs := make([]uintptr, 1)
	n := runtime.Callers(5, fpcs)
	relpath, name, row, err := getCaller(n, fpcs)

	if maxDepthOfCallerPath > 0 {
		pathElements := strings.Split(relpath, string(os.PathSeparator))
		length := len(pathElements)

		start := length - maxDepthOfCallerPath
		if start < 0 {
			start = 0
		} else {
			relpath = "..." + string(os.PathSeparator)
		}
		relpath += path.Join(pathElements[start:length]...)
	}

	if err == nil {
		caller.Path = relpath
		caller.FunctionName = name
		caller.LineNumber = row
	}

	msg := Message{
		Time:    now,
		Level:   level,
		Caller:  caller,
		Message: fmt.Sprint(params...),
	}

	return msg
}

// logHandler calls all defined handlers with the built Message object
func logHandler(level LogLevel, params ...interface{}) {
	if !showMe(level) {
		return
	}

	if config == nil {
		config = DefaultLevelConfig()
	}

	lvlName := strings.Title(strings.ToLower(level.String()))
	s := structs.New(config)
	lvlField := s.Field(lvlName)
	cfg := lvlField.Value().(LevelConfig)

	message := buildMessage(level, params...)

	for _, handler := range cfg.Handlers {
		handler(message)
	}

}

// log is the internal log handler
func log(message Message) {

	logMessage := stringify(message)

	fmt.Print(logMessage)
}

func println(level LogLevel, params ...interface{}) {
	params = append(params, "\n")
	logHandler(level, params...)
}

func sprint(level LogLevel, params ...interface{}) string {
	if !showMe(level) {
		return ""
	}
	message := buildMessage(level, params...)
	return stringify(message)
}

func sprintln(level LogLevel, params ...interface{}) string {
	params = append(params, "\n")
	return sprint(level, params...)
}

func showMe(level LogLevel) bool {
	if logLevel == NONE || level == NONE {
		return false
	}

	return logLevel >= level
}

func getLogLevel(withFormat bool, values ...interface{}) (logLevel LogLevel, format string, newValues []interface{}) {

	comp := LogLevel(0)
	try := values[0]
	level := defaultLevel

	if reflect.TypeOf(try) == reflect.TypeOf(comp) {
		val := reflect.ValueOf(try)
		level = LogLevel(val.Uint())
		values = values[1:]
	}

	format = ""
	if withFormat {
		try := values[0]

		if reflect.TypeOf(try) == reflect.TypeOf(format) {
			val := reflect.ValueOf(try)
			format = val.String()
			values = values[1:]
		} else {
			log2.Panicln("please specify a format")
		}
	}

	return level, format, values
}

func getCaller(n int, fpcs []uintptr) (relpath string, name string, row int, err error) {
	err = nil

	if n == 0 {
		return "", "", -1, errors.New("MSG CALLER WAS NIL")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return "", "", -1, errors.New("MSG CALLER WAS NIL")
	}

	// Get Path
	absPath, row := caller.FileLine(fpcs[0] - 1)
	dir, err := os.Getwd()
	if err != nil {
		log2.Fatal(err)
	}
	relpath, err = filepath.Rel(dir, absPath)
	if err != nil {
		return "", "", -1, err
	}

	// Get Name of the caller function
	na := strings.Split(caller.Name(), ".")
	name = na[len(na)-1]

	return
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
