# AwesomeLog

AwsomeLog is a inplace replacement of the default log package with some extended features. 

Mainly **AwsomeLog** provides functionality to define log levels as well as a **PrettyPrint** function to niceley and readable print out objects. 
It also shows the filename and line of code on logs  with a level of **DEBUG** or **VERBOSE**


### Quick start
Import the module with the alias `log` as shown in the example below.

To get log output you **MUST** define what log levels should be printed out. 
**IMPORTANT** The default log level is `NONE`.

```go
package main
import (
    log "github.com/chris-dot-exe/AwsomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)
}
``` 

### Functions

The following functions are provided from this package:
```go
log.Print()
log.Println()
log.Printf()
log.PrettyPrint()
```

Functions without additional functionality. Panic and Fatal directly calls the original log functions:
```go
log.Fatal()
log.Fatalf()
log.Fatalln()
log.Panic()
log.Panicf()
log.Panicln()
```

And AwesomeLog setting functions:
```go
log.SetLogLevel()
log.SetLogLEvelByString()
log.SetDefaultLevel()
log.ShowColorsInLogs()
```

### Setup Functions:
#### `log.SetLogLevel(logLevel)`
`log.SetLogLevel()` defines to which level messages should be logged. 

#### `log.SetLogLevelByString(string)`
`log.SetLogLevelByString()` same as `SetLogLevel()` but you can pass the log level as a string to the function (e.g. from a configuration file)

**Default is `log.NONE` / `NONE`** 


See examples below.

#### `log.SetDefaultLevel(logLevel)`
`log.SetDefaultLevel()` defines the default log level if a print function is called without a log level as first argument. 

Default: log.INFO

#### `log.ShowColorsInLogs(bool)`
`log.ShowColorsInLogs()` defines if the colored log-level labels should be shown in logs which are redirected to a file. 
If this is set to false and the outout is visible in the terminal AND is saved to a log file (e.g. docker logs) it will be shown with colors on the terminal output but without in the docker log.  

**Default: false**

This function is not fully tested. 



## Examples
```go
package main
import (
    log "github.com/chris-dot-exe/AwsomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Print("foobar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwsomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Println("Foobar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwsomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Printf("Foo%s", "Bar")
}
```

```go
package main
import (
    log "github.com/chris-dot-exe/AwsomeLog"
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
    log "github.com/chris-dot-exe/AwsomeLog"
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
    log "github.com/chris-dot-exe/AwsomeLog"
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

## Notes
PrettyPrint can only display fields which are exported!
