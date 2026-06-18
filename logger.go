//go:generate go run generate.go
package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	log2 "log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

func init() {
	config = DefaultLevelConfig()
	isTerminal = checkIsTerminal()
	dir, err := os.Getwd()
	if err != nil {
		log2.Fatalf("Failed to get current working directory: %v", err)
	}
	rootedPath = dir
}

// SetLogLevel defines to which LogLevel log messages should be shown.
//
// Default is VERBOSE
func SetLogLevel(lvl LogLevel) {
	mu.Lock()
	defer mu.Unlock()

	logLevel = lvl
}

func setCallerForLevel(levelCfg *LevelConfig, show bool) {
	levelCfg.ShowFunctionName = show
	levelCfg.ShowFilePath = show
	levelCfg.ShowLineNumber = show
}

// ShowCaller defines if the caller (function name, line number, file path) should be shown on a global level.
func ShowCaller(show bool) {
	mu.Lock()
	defer mu.Unlock()
	if config == nil {
		config = DefaultLevelConfig()
	}
	// Nutzt die vom Generator erstellte Funktion
	applyShowCaller(config, show)
}

// ShowColors Defines if colored level tags should be shown in the console log.
func ShowColors(show bool) {
	mu.Lock()
	defer mu.Unlock()
	showColors = show
}

// SetLogLevelByString defines to which LogLevel log messages should be shown based on the given string e.g. SetLogLevelByString("WARN")
// This is useful if the LogLevel is defined in a config file.
func SetLogLevelByString(lvlStr string) {
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	defer mu.Unlock()
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
	mu.Lock()
	defer mu.Unlock()
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
		Info: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
			Handlers:         []Handler{log},
		},
		Warn: LevelConfig{
			ShowLineNumber:   false,
			ShowFunctionName: false,
			ShowFilePath:     false,
			Handlers:         []Handler{log},
		},
		Error: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
			Handlers:         []Handler{log},
		},
		Critical: LevelConfig{
			ShowLineNumber:   true,
			ShowFunctionName: true,
			ShowFilePath:     true,
			Handlers:         []Handler{log},
		},
	}
	return cfg
}

// SetLevelConfig set the config for AwesomeLog
func SetLevelConfig(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	config = cfg
}

func SetCallerSkip(skip int) {
	mu.Lock()
	defer mu.Unlock()
	callerSkip = skip
}

// SetTimeFormat set the timeformat for log messages
func SetTimeFormat(format string) {
	mu.Lock()
	defer mu.Unlock()
	timeFormat = format
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	out = w
	if w == nil {
		isTerminal = isatty.IsTerminal(os.Stdout.Fd())
	} else if f, ok := w.(*os.File); ok {
		isTerminal = isatty.IsTerminal(f.Fd())
	} else {
		isTerminal = false
	}
}

// Println logs a message at the defined LogLevel a newline is appended
func Println(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	println(level, params...)
}

// Print logs a message at the defined LogLevel
func Print(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	print(level, params...)
}

// Printf logs a message at the defined LogLevel and formats the message according to a format specifier
func Printf(paramsOriginal ...interface{}) {
	level, format, params := getLogLevel(true, paramsOriginal...)
	print(level, fmt.Sprintf(format, params...))
}

// PrettyPrint logs a message at the defined LogLevel formatted as JSON
// Works only with exported fields.
func PrettyPrint(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)

	for i, param := range params {
		if masker, ok := param.(LogMasker); ok {
			params[i] = masker.LogValue()
		}
	}

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

	for i, param := range params {
		if masker, ok := param.(LogMasker); ok {
			params[i] = masker.LogValue()
		}
	}

	b, err := json.MarshalIndent(params, "", "  ")
	if err != nil {
		Fatal("unsupported input. error: ", err)
	}

	return sprintln(level, string(b))
}

// region fatal

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

// endregion fatal

// region panic

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

// endregion panic

// stringify builds the log message string with colors and caller
func stringify(message Message) string {

	if config == nil {
		config = DefaultLevelConfig()
	}

	cfg := getConfigForLevel(message.Level)

	var sb strings.Builder
	sb.Grow(80 + len(message.Message))

	//prefix := ""
	//caller := ""

	if showTimestamp {
		//prefix = fmt.Sprintf("%s ", message.Time.Format(timeFormat))
		sb.WriteString(message.Time.Format(timeFormat))
		sb.WriteString(" ")
	}

	if showColors && (colorsInLogs || isTerminal) {
		//prefix += fmt.Sprintf(message.Level.Color()+"[%s]"+ANSI_RESET, message.Level.String())
		sb.WriteString(message.Level.Color())
		sb.WriteString("[")
		sb.WriteString(message.Level.String())
		sb.WriteString("]")
		sb.WriteString(ANSI_RESET)
	} else {
		//prefix += fmt.Sprintf("[%s]", message.Level.String())
		sb.WriteString("[")
		sb.WriteString(message.Level.String())
		sb.WriteString("]")
	}

	if cfg.ShowFilePath || cfg.ShowFunctionName || cfg.ShowLineNumber {
		//caller += "["
		sb.WriteString("[")
		if cfg.ShowFilePath {

			//caller += fmt.Sprintf("%s:", message.Caller.Path)
			sb.WriteString(message.Caller.Path)
			sb.WriteString(":")
		}
		if cfg.ShowFunctionName {
			//caller += fmt.Sprintf("%s", message.Caller.FunctionName)
			sb.WriteString(message.Caller.FunctionName)
		}
		if cfg.ShowLineNumber {
			//caller += fmt.Sprintf(":%d", message.Caller.LineNumber)
			sb.WriteString(":")
			sb.WriteString(fmt.Sprint(message.Caller.LineNumber))
		}

		//caller += "]"
		sb.WriteString("]")
	}

	sb.WriteString(" ")
	sb.WriteString(message.Message)

	//return fmt.Sprintf("%s%s %s", prefix, caller, message.Message)
	return sb.String()
}

// buildMessage builds the Message object used by all log handlers
func buildMessage(level LogLevel, params ...interface{}) Message {
	now := time.Now()
	caller := Caller{}

	fpcs := make([]uintptr, 1)
	n := runtime.Callers(callerSkip, fpcs)
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

	cfg := getConfigForLevel(level)
	message := buildMessage(level, params...)

	for _, handler := range cfg.Handlers {
		handler(message)
	}

}

// log is the internal log handler
func log(message Message) {

	logMessage := stringify(message)

	mu.RLock()
	defer mu.RUnlock()
	w := out
	if w == nil {
		w = os.Stdout
	}
	fmt.Fprint(w, logMessage)
}

func println(level LogLevel, params ...interface{}) {
	params = append(params, "\n")
	logHandler(level, params...)
}

func print(level LogLevel, params ...interface{}) {
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
	mu.RLock()
	defer mu.RUnlock()

	if logLevel == NONE || level == NONE {
		return false
	}

	return logLevel >= level
}

func getLogLevel(withFormat bool, values ...interface{}) (logLevel LogLevel, format string, newValues []interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	if len(values) == 0 {
		values = append(values, "")
	}

	level := defaultLevel

	if lvl, ok := values[0].(LogLevel); ok {
		level = lvl
		values = values[1:]
	}

	format = ""
	if withFormat {
		try := values[0]

		if f, ok := try.(string); ok {
			format = f
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
	//relpath, err = filepath.Rel(rootedPath, absPath)
	//if err != nil {
	//	return "", "", -1, err
	//}
	//
	//// Get Name of the caller function
	//na := strings.Split(caller.Name(), ".")
	//name = na[len(na)-1]

	relpath = absPath
	if strings.HasPrefix(relpath, rootedPath) {
		relpath = absPath[len(rootedPath):]
		if len(relpath) > 0 && (relpath[0] == '/' || relpath[0] == '\\') {
			relpath = relpath[1:]
		}
	}

	if maxDepthOfCallerPath > 0 {
		slashesFound := 0
		for i := len(relpath) - 1; i >= 0; i-- {
			if relpath[i] == '/' || relpath[i] == '\\' {
				slashesFound++
			}
			if slashesFound == maxDepthOfCallerPath {
				relpath = relpath[i+1:]
				break
			}
		}
	}

	fullName := caller.Name()
	name = fullName
	if lastSlash := strings.LastIndexByte(fullName, '/'); lastSlash >= 0 {
		name = fullName[lastSlash+1:]
	}
	if lastDot := strings.LastIndexByte(name, '.'); lastDot >= 0 {
		name = name[lastDot+1:]
	}

	return relpath, name, row, nil
}

func checkIsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}
