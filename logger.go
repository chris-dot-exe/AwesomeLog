package log

import (
	"encoding/json"
	"errors"
	"fmt"
	log2 "log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/fatih/structs"
	"golang.org/x/crypto/ssh/terminal"
)

type LogLevel uint

func (t *LogLevel) String() string {
	for k, v := range level {
		if *t == v {
			return k
		}
	}
	return "NONE"
}

func (t *LogLevel) Color() string {
	return lvlColor[*t]
}

const (
	NONE    LogLevel = 0
	WARN    LogLevel = 1
	INFO    LogLevel = 2
	DEBUG   LogLevel = 10
	VERBOSE LogLevel = 20
)

type LevelConfig struct {
	ShowLineNumber   bool
	ShowFunctionName bool
	ShowFilePath     bool
}

type Config struct {
	Verbose LevelConfig
	Debug   LevelConfig
	Warn    LevelConfig
	Info    LevelConfig
}

var logLevel = VERBOSE
var defaultLevel = INFO
var colorsInLogs = false
var showColors = true
var config *Config
var showTimestamp = true

var level = map[string]LogLevel{
	"NONE":    NONE,
	"WARN":    WARN,
	"INFO":    INFO,
	"DEBUG":   DEBUG,
	"VERBOSE": VERBOSE,
}

var lvlColor = map[LogLevel]string{
	NONE:    "",
	WARN:    ANSI_YELLOW_BACKGROUND + ANSI_BLACK,
	INFO:    ANSI_BLUE_BACKGROUND + ANSI_WHITE,
	DEBUG:   ANSI_RED_BACKGROUND + ANSI_WHITE,
	VERBOSE: "",
}

func SetLogLevel(lvl LogLevel) {
	logLevel = lvl
}

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

func ShowColors(show bool) {
	showColors = show
}

func SetLogLevelByString(lvlStr string) {
	lvlStr = strings.ToUpper(lvlStr)
	val, ok := level[lvlStr]
	if !ok {
		log2.Fatalf("LogLevel '%s' is not supported!\n", lvlStr)
		return
	}
	logLevel = val
}

func SetDefaultLevel(lvl LogLevel) {
	defaultLevel = lvl
}

func ShowColorsInLogs(show bool) {
	colorsInLogs = show
}

func ShowTimestamp(show bool) {
	showTimestamp = show
}

func DefaultLevelConfig() *Config {
	cfg := &Config{
		Verbose: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
		},
		Debug: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
		},
		Warn: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
		},
		Info: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
		},
	}
	return cfg
}

func SetLevelConfig(cfg *Config) {
	config = cfg
}

func Println(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	println(level, params...)
}

func Print(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	print(level, params)
}

func Printf(paramsOriginal ...interface{}) {
	level, format, params := getLogLevel(true, paramsOriginal...)
	print(level, fmt.Sprintf(format, params...))
}

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

func Fatal(params ...interface{}) {
	log2.Fatal(params...)
}

func Fatalf(format string, params ...interface{}) {
	log2.Fatalf(format, params...)
}

func Fatalln(params ...interface{}) {
	log2.Fatalln(params...)
}

func Panic(params ...interface{}) {
	log2.Panic(params...)
}

func Panicf(format string, params ...interface{}) {
	log2.Panicf(format, params...)
}

func Panicln(params ...interface{}) {
	log2.Panicln(params...)
}

func stringify(level LogLevel, params ...interface{}) string {
	if config == nil {
		config = DefaultLevelConfig()
	}

	lvlName := strings.Title(strings.ToLower(level.String()))

	s := structs.New(config)
	lvlField := s.Field(lvlName)

	cfg := lvlField.Value().(LevelConfig)

	prefix := ""
	caller := ""
	if showColors && (colorsInLogs || isTerminal()) {
		prefix = fmt.Sprintf(level.Color()+"[%s]"+ANSI_RESET, level.String())
	} else {
		prefix = fmt.Sprintf("[%s]", level.String())
	}
	if cfg.ShowFilePath || cfg.ShowFunctionName || cfg.ShowLineNumber {
		caller += "["
		fpcs := make([]uintptr, 1)
		n := runtime.Callers(5, fpcs)
		relpath, name, row, err := getCaller(n, fpcs)
		if err == nil {
			if cfg.ShowFilePath {
				caller += fmt.Sprintf("%s:", relpath)
			}
			if cfg.ShowFunctionName {
				caller += fmt.Sprintf("%s", name)
			}
			if cfg.ShowLineNumber {
				caller += fmt.Sprintf(":%d", row)
			}

		}
		caller += "]"
	}

	return fmt.Sprintf("%s%s %s", prefix, caller, fmt.Sprint(params...))
}

func print(level LogLevel, params ...interface{}) {
	if !showMe(level) {
		return
	}
	if !showTimestamp {
		fmt.Print(stringify(level, params...))
		return
	}

	log2.Print(stringify(level, params...))
}

func println(level LogLevel, params ...interface{}) {
	params = append(params, "\n")
	print(level, params...)
}

func sprint(level LogLevel, params ...interface{}) string {
	if !showMe(level) {
		return ""
	}
	return stringify(level, params...)
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
