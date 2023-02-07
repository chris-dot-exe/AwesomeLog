AwesomeLog
=========
[![Go Reference](https://pkg.go.dev/badge/github.com/chris-dot-exe/AwesomeLog.svg)](https://pkg.go.dev/github.com/chris-dot-exe/AwesomeLog)
![Last Release Version](https://img.shields.io/github/v/release/chris-dot-exe/AwesomeLog?include_prereleases)
![Go Version](https://img.shields.io/github/go-mod/go-version/chris-dot-exe/AwesomeLog/master)
[![License](https://img.shields.io/github/license/chris-dot-exe/AwesomeLog)](https://github.com/chris-dot-exe/AwesomeLog/blob/master/LICENSE)

AwesomeLog is a fully compatible drop-in replacement for the standard library logger with some awesome features.

AwesomeLog let you define log levels for each logged messages as well as a **PrettyPrint** function to print out objects in a pretty readable format.
It also adds the option to show details of the _caller_ like file path, function name and line number.

AwesomeLog now also provides the functionality to add custom handlers for each LogLevel.

### Documentation
[Documentation](https://pkg.go.dev/github.com/chris-dot-exe/AwesomeLog@v1.0.0-rc#section-documentation)

### Quick start
The simplest way to use AwesomeLog is to just replace the standard library log with AwesomeLog:

```go
package main

import (
log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
log.Println("Hello World!")
}
```
Output:
`2022/02/14 09:32:44 [INFO] Hello World!`


## Examples
```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Print("foobar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Println("Foobar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Printf("Foo%s", "Bar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

type foo struct {
	Foo string
	Bar string
	Foobar struct {
		Meeps []string
	}
}

func main() {
    log.SetLogLevel(log.VERBOSE)

    foo := foo{
		Foo: "Test",
		Bar: "Test",
		Foobar: struct {
    	Meeps []string
		}{[]string{"Meep", "Meep2", "Meep2.1"}},
	}

    log.PrettyPrint(foo)
}
```

Output:

<img alt="cmdline output" src="https://user-images.githubusercontent.com/49272981/80649110-b3297e80-8a71-11ea-9779-d359da872d75.png" width="500px">



### Examples with Loglevel
Now the interesting part:

```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Println(log.VERBOSE, "Foobar Verbose")
    log.Println(log.DEBUG, "Foobar Debug")
    log.Println(log.INFO, "Foobar Info")
    log.Println(log.WARN, "Foobar Warning")
}
```
Output:

<img alt="cmdline output" src="https://user-images.githubusercontent.com/49272981/80649108-b290e800-8a71-11ea-8463-595e9de9f171.png" width="500px">

### Show only messages to a specific level:
The priority of the log levels is as following (highest to lowest):

```
NONE
WARN
INFO
DEBUG
VERBOSE
```

If you set the log-level to info:
```go
package main
import (
    log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.INFO)

    log.Println(log.VERBOSE, "Foobar Verbose")
    log.Println(log.DEBUG, "Foobar Debug")
    log.Println(log.INFO, "Foobar Info")
    log.Println(log.WARN, "Foobar Warning")
}
```
The output is reduced to the following messages:

<img alt="cmdline output" src="https://user-images.githubusercontent.com/49272981/80649104-b1f85180-8a71-11ea-9c72-d98ed825a6b4.png" width="500px">

### Config Example

```go
package main
import (
	log "github.com/chris-dot-exe/AwesomeLog"
)

func main() {
	cfg := log.DefaultLevelConfig()

    cfg.Debug.ShowLineNumber = false
    cfg.Debug.ShowFunctionName = true
    cfg.Debug.ShowFilePath = false

	log.SetLevelConfig(cfg)
}
```

### Custom Handler
It is possible to add custom handler for each LogLevel.<br>
The example below shows how a custom handler for GlitchTip/Sentry can be defined: 
```go
package main

import (
  log "github.com/chris-dot-exe/AwesomeLog"
  "github.com/getsentry/sentry-go"
  "time"
)

func main() {
  // Setup GlitchTip
  sentry.Init(sentry.ClientOptions{
    Dsn: "http://0d985cf763a34732a4839eea121c2f25@localhost:8000/1",
  })
  defer sentry.Flush(time.Second * 5)
  // Setup AwesomeLog
  log.SetDefaultLevel(log.INFO)

  // Get Default Level Config
  lvlConfig := log.DefaultLevelConfig()
  // Add custom Handler
  lvlConfig.Warn.AddHandler(GlitchContextLogger)
  lvlConfig.Debug.AddHandler(GlitchContextLogger)
  lvlConfig.Info.AddHandler(GlitchMessage)
  // Set new Level Config
  log.SetLevelConfig(lvlConfig)

  logTest()
}

func GlitchContextLogger(message log.Message) {
  sentry.ConfigureScope(func(scope *sentry.Scope) {
    scope.SetContext("caller", message.Caller)
    scope.SetTag("level", message.Level.String())
  })
  sentry.CaptureMessage(message.Message)
}

func GlitchMessage(message log.Message) {
  sentry.CaptureMessage(message.Message)
}

func logTest() {
  log.Println(log.DEBUG, "some debug message")
  log.Println(log.INFO, "some info message")
  log.Println(log.WARN, "something went wrong")
  log.Println(log.VERBOSE, "verbose message, not sent to GlitchTip")
}
```
In GlitchTip:

<a href="https://user-images.githubusercontent.com/49272981/154164691-e95b7005-71b8-4b1d-8e8e-b172be3f5de7.png">
<img alt="GlitchTip Details" src="https://user-images.githubusercontent.com/49272981/154164691-e95b7005-71b8-4b1d-8e8e-b172be3f5de7.png" width="500px">
</a>
<a href="https://user-images.githubusercontent.com/49272981/154164688-0f21ce4a-1140-47cc-abf1-39903688a782.png">
<img alt="GlitchTip Details" src="https://user-images.githubusercontent.com/49272981/154164688-0f21ce4a-1140-47cc-abf1-39903688a782.png" width="500px">
</a>
