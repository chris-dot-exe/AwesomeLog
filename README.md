AwesomeLog
=========
[![Go Reference](https://pkg.go.dev/badge/github.com/chris-dot-exe/AwesomeLog.svg)](https://pkg.go.dev/github.com/chris-dot-exe/AwesomeLog)
![Last Release Version](https://img.shields.io/github/v/release/chris-dot-exe/AwesomeLog?include_prereleases)
![Go Version](https://img.shields.io/github/go-mod/go-version/chris-dot-exe/AwesomeLog/master)
[![License](https://img.shields.io/github/license/chris-dot-exe/AwesomeLog)](https://github.com/chris-dot-exe/AwesomeLog/blob/master/LICENSE)


# AwesomeLog
AwesomeLog is a lightweight, blazing fast, and thread-safe logger for Go that can be used as a direct drop-in replacement for the standard log package.

With this major performance update, the logger has been optimized from the ground up to maximize speed (near zero-allocations), extensibility, and convenience, while retaining full backward compatibility.

Logging is now up to nine times faster than before and uses up to 10 times fewer allocations.

## ✨ New Features
* **Extreme Performance:** By eliminating fmt.Sprintf at its core and avoiding expensive reflection (wherever possible), AwesomeLog is now blazing fast.
* **Code Generation:** All level-specific functions (Error(), Infof(), etc.) are generated and call the internal logging engine directly, removing any dynamic overhead.
* **Data Masking (LogMasker):** You can easily mask sensitive data (passwords, tokens) in your structs before logging them as JSON, WITHOUT affecting your REST API representations.
* **Custom Output Streams:** Seamlessly redirect logs to files, network sockets, or buffers. AwesomeLog automatically detects whether the new output supports ANSI colors and adapts accordingly.

## Performance Comparison

```
goos: linux
goarch: amd64
pkg: github.com/chris-dot-exe/AwesomeLog
cpu: AMD Ryzen 7 7800X3D 8-Core Processor           
BenchmarkAwesomeLog_Info-16              1714042               680.8 ns/op           152 B/op          3 allocs/op
BenchmarkAwesomeLog_Println-16           1606201               740.0 ns/op           184 B/op          4 allocs/op
BenchmarkStdlibLog_Println-16           14904094               80.34 ns/op            24 B/op          1 allocs/op
BenchmarkAwesomeLog_Info_Full-16         1254553               953.1 ns/op           386 B/op          6 allocs/op
BenchmarkAwesomeLog_Println_Full-16      1000000                1030 ns/op           418 B/op          7 allocs/op
BenchmarkStdlibLog_Println_Full-16       1965724               612.8 ns/op           256 B/op          3 allocs/op
PASS
ok      github.com/chris-dot-exe/AwesomeLog     10.134s

#Before: 
goos: linux
goarch: amd64
pkg: github.com/chris-dot-exe/AwesomeLog
cpu: AMD Ryzen 7 7800X3D 8-Core Processor           
BenchmarkAwesomeLog_Info-16               187358              6201 ns/op            1737 B/op         30 allocs/op
BenchmarkAwesomeLog_Println-16            185704              6252 ns/op            1777 B/op         31 allocs/op
BenchmarkStdlibLog_Println-16            2258814             529.0 ns/op              24 B/op          1 allocs/op
BenchmarkAwesomeLog_Info_Full-16          173394              6799 ns/op            2249 B/op         44 allocs/op
BenchmarkAwesomeLog_Println_Full-16       173806              6854 ns/op            2281 B/op         45 allocs/op
BenchmarkStdlibLog_Println_Full-16        871966              1214 ns/op             256 B/op          3 allocs/op
PASS
ok      github.com/chris-dot-exe/AwesomeLog     7.776s
```


### Documentation
[Documentation](https://pkg.go.dev/github.com/chris-dot-exe/AwesomeLog@v1.0.0-rc#section-documentation)

### Installation

```sh
go get github.com/chris-dot-exe/AwesomeLog
```

## 📖 Features & Examples
1. Convenience Methods
   Instead of passing the log level as a parameter every time, you can now use dedicated level functions directly:

```go
package main
import "github.com/chris-dot-exe/AwesomeLog"
func main() {
  log.Info("This is an info message!")
  log.Warnf("An warning occurred: %s\n", err)
  log.Errorfln("A error occurred: %s", err)
  log.Debug("Debugging enabled")
}
```

2. Custom Output (SetOutput)
   Want to write your logs to a file? No problem. AwesomeLog automatically disables ANSI colors if the output is not a terminal.

```go
file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log.SetOutput(file)
log.Info("This will be written to the file without ANSI colors")
```

3. Data Masking (LogMasker)
   If you use log.PrettyPrint() to output objects as JSON, you might want to hide passwords. Simply implement the LogMasker interface!

```go
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
// LogValue is called exclusively by AwesomeLog!
func (u User) LogValue() interface{} {
    return struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }{
        Username: u.Username,
        Password: "***REDACTED***",
    }
}
func main() {
    u := User{Username: "admin", Password: "supersecret"}
    log.PrettyPrint(u) // Password will be logged as "***REDACTED***"
}
```

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
    log.Println(log.ERROR, "Foobar Error")
    log.Println(log.CRITICAL, "Foobar PANIC MODE! aka critical")
}
```
Output:

<img alt="cmdline output" src="https://user-images.githubusercontent.com/49272981/247899663-c83072a5-e6d8-420c-8dda-2c3b9dca6916.png" width="650px">

### Show only messages to a specific level:
The priority of the log levels is as following (highest to lowest):

```
NONE
CRITICAL
ERROR
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
    log.Println(log.ERROR, "Foobar Error")
    log.Println(log.CRITICAL, "Foobar Critical")
}
```
The output is reduced to the following messages:

<img alt="cmdline output" src="https://user-images.githubusercontent.com/49272981/247900382-9f02cf3a-51bd-4c75-a82f-bfa25f8ceade.png" width="650px">

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
