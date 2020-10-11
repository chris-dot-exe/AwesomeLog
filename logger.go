package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	log2 "log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type LogLevel uint

const (
	NONE LogLevel = 0
	WARN LogLevel = 1
	INFO LogLevel = 2
	DEBUG LogLevel = 10
	VERBOSE LogLevel = 20
)

var LOG_LEVEL = VERBOSE
var DEFAULT_LEVEL = INFO
var COLORS_IN_LOGS = false

var level = map[string]LogLevel{
	"NONE": NONE,
	"WARN": WARN,
	"INFO": INFO,
	"DEBUG": DEBUG,
	"VERBOSE": VERBOSE,
}

var lvlColor = map[LogLevel]string{
	NONE: "",
	WARN: ANSI_YELLOW_BACKGROUND+ANSI_BLACK,
	INFO: ANSI_BLUE_BACKGROUND+ANSI_WHITE,
	DEBUG: ANSI_RED_BACKGROUND+ANSI_WHITE,
	VERBOSE: "",
}


func SetLogLevel(lvl LogLevel) {
	LOG_LEVEL = lvl
}

func SetLogLevelByString(lvlStr string) {
	val, ok := level[lvlStr]
	if !ok {
		log2.Fatalf("LogLevel '%s' is not supported!\n", lvlStr)
		return
	}
	LOG_LEVEL = val
}

func SetDefaultLevel(lvl LogLevel)  {
	DEFAULT_LEVEL = lvl
}

func ShowColorsInLogs(ok bool) {
	COLORS_IN_LOGS = ok
}

func Println(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)

	if !showMe(level) {
		return
	}
	fpcs := make([]uintptr, 1)
	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)

	_, name, row, err := getCaller(n, fpcs)
	if err != nil || level < 10  {
		log2.Println(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +" %s", getLevelName(level), fmt.Sprint(params...)))
		return
	}

	log2.Println(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +"[%s:%d] %s", getLevelName(level), name, row, fmt.Sprint(params...)))

}

func getLogLevel(withFormat bool, values ...interface{}) (logLevel LogLevel, format string, newValues []interface{}) {

	comp := LogLevel(0)
	try := values[0]
	level := DEFAULT_LEVEL

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

func PrettyPrint( params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	if !showMe(level) {
		return
	}
	fpcs := make([]uintptr, 1)
	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)

	_, name, row, err := getCaller(n, fpcs)

	b, err := json.MarshalIndent(params, "", "  ")
	if err != nil || level < 10 {
		if COLORS_IN_LOGS || isTerminal() {
			log2.Println(fmt.Sprintf(getLevelColor(level)+"[%s]"+ANSI_RESET+"\n%s", getLevelName(level), string(b)))
			return
		}
		log2.Println(fmt.Sprintf("[%s]\n%s", getLevelName(level), string(b)))
		return
	}
	if COLORS_IN_LOGS || isTerminal() {
		log2.Println(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +"[%s:%d]\n%s", getLevelName(level), name, row, string(b)))
		return
	}
	log2.Println(fmt.Sprintf("[%s][%s:%d]\n%s", getLevelName(level), name, row, string(b)))
	return


}

func Print(params ...interface{}) {
	level, _, params := getLogLevel(false, params...)
	if !showMe(level) {
		return
	}
	fpcs := make([]uintptr, 1)
	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)

	_, name, row, err := getCaller(n, fpcs)
	if err != nil || level < 10 {
		if COLORS_IN_LOGS || isTerminal() {
			log2.Print(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +" %s", getLevelName(level), fmt.Sprint(params...)))
			return
		}
		log2.Print(fmt.Sprintf("[%s] %s", getLevelName(level), fmt.Sprint(params...)))
		return
	}
	if COLORS_IN_LOGS || isTerminal() {
		log2.Print(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +"[%s:%d] %s", getLevelName(level), name, row, fmt.Sprint(params...)))
		return
	}
	log2.Print(fmt.Sprintf("[%s][%s:%d] %s", getLevelName(level), name, row, fmt.Sprint(params...)))
	return

}

func Printf(paramsOriginal ...interface{}) {
	level, format, params := getLogLevel(true, paramsOriginal...)
	if !showMe(level) {
		return
	}
	fpcs := make([]uintptr, 1)
	// Skip 2 levels to get the caller
	n := runtime.Callers(2, fpcs)

	_, name, row, err := getCaller(n, fpcs)
	if err != nil || level < 10 {
		if COLORS_IN_LOGS || isTerminal() {
			log2.Printf(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +" ", getLevelName(level)) + format, params...)
			return
		}
		log2.Printf(fmt.Sprintf("[%s] ", getLevelName(level)) + format, params...)
		return
	}
	if COLORS_IN_LOGS || isTerminal() {
		log2.Printf(fmt.Sprintf(getLevelColor(level)+"[%s]" + ANSI_RESET +"[%s:%d] ", getLevelName(level), name, row,) + format, params...)
		return
	}
	log2.Printf(fmt.Sprintf("[%s][%s:%d] ", getLevelName(level), name, row,) + format, params...)
	return

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



func showMe(level LogLevel) bool {
	if LOG_LEVEL == NONE {
		return false
	} else if LOG_LEVEL == VERBOSE {
		return true
	} else if LOG_LEVEL == DEBUG && level <= DEBUG  {
		return true
	} else if LOG_LEVEL == INFO && level <= INFO  {
		return true
	} else if LOG_LEVEL == WARN && level <= WARN {
		return true
	} else {
		return false
	}
}



func getCaller(n int, fpcs []uintptr) (relpath string, name string, row int, err error) {
	err = nil

	if n == 0 {
		return "", "", -1, errors.New("MSG CALLER WAS NIL")
	}

	caller := runtime.FuncForPC(fpcs[0]-1)
	if caller == nil {
		return "", "", -1, errors.New("MSG CALLER WAS NIL")
	}

	// Get Path
	absPath, row := caller.FileLine(fpcs[0]-1)
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

func getLevelName(lvl LogLevel) string {
	for k, v := range level {
		if lvl == v {
			return k
		}
	}
	return "NONE"
}

func getLevelColor(lvl LogLevel) string {
	return lvlColor[lvl]
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}