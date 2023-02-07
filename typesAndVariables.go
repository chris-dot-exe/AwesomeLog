package log

type LogLevel uint

// String returns the LogLevel name as string
func (t *LogLevel) String() string {
	for k, v := range level {
		if *t == v {
			return k
		}
	}
	return "NONE"
}

// Color returns the defined color of the LogLevel
func (t *LogLevel) Color() string {
	return lvlColor[*t]
}

// Caller contains the Caller information file path, function name and line number
type Caller struct {
	Path         string
	FunctionName string
	LineNumber   int
}

// Message is the object which is passed to every handler function.
// Message contains the LogLevel, the Caller object and the message
type Message struct {
	Level   LogLevel
	Caller  Caller
	Message string
}

type Handler func(message Message)

// LevelConfig represents the configuration for each LogLevel
type LevelConfig struct {
	ShowLineNumber   bool
	ShowFunctionName bool
	ShowFilePath     bool
	Handlers         []Handler
}

// AddHandler adds a custom Handler to the existing handlers of the LogLevel
func (c *LevelConfig) AddHandler(handler Handler) {
	c.Handlers = append(c.Handlers, handler)
}

// SetHandlers sets custom handlers for the LogLevel.
// SetHandlers overrides the existing handler
func (c *LevelConfig) SetHandlers(handler []Handler) {
	c.Handlers = handler
}

// Config represents the config for all LogLevels
type Config struct {
	Verbose LevelConfig
	Debug   LevelConfig
	Warn    LevelConfig
	Info    LevelConfig
}

const (
	NONE    LogLevel = 0
	WARN    LogLevel = 1
	INFO    LogLevel = 2
	DEBUG   LogLevel = 10
	VERBOSE LogLevel = 20
)

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
