# AwesomeLog

AwesomeLog ist eine inplace Erweiterung des Standard log packages von @chris.exe
Welche bei der Entwicklung eines Go Projekts entstand. 

AwesomeLog bietet die Möglichkeit verschiedene Log Level zu definieren.
Sowie eine Funktion die Objekte schöner darstellt.

### Quick start
Das Module wird mit dem Alias log importiert und ersetzt so die Standard Log Funktion.

Mun eine Log Ausgabe zu erhalten muss nun lediglich definiert werden bis zu welchem Loglevel Logs ausgegeben werden (Default: None).

```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)
}
``` 

### Funktionen

Implementierte Funktionen sind `log.Println()`, `log.Print()`, `log.Printf()` und `log.PrettyPrint()`


## Examples
```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Print("foobar")
}
```

```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Println("Foobar")
}
```

```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Printf("Foo%s", "Bar")
}
```

```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
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

Ausgaben:

![cmdline output](https://intern.t3debug.de/git/sg-go-libs/AwesomeLog/-/wikis/uploads/c0ec95041f009a18f7d4e045a763cad3/Auswahl_305.png)


###Examples with Loglevel
So kommen wir jetzt zum interessanten Teil!

```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.VERBOSE)

    log.Println(log.VERBOSE, "Foobar Verbose")
    log.Println(log.DEBUG, "Foobar Debug")
    log.Println(log.INFO, "Foobar Info")
    log.Println(log.WARN, "Foobar Warning")
}
```
Ausgaben:

![cmdline output](https://intern.t3debug.de/git/sg-go-libs/AwesomeLog/-/wikis/uploads/99606d30fc88e6fe045fbfdee7255ec6/Auswahl_306.png)

### Nur bestimmte Loglevel anzeigen:
Will man nun, nach dem testen, die Logausgabe reduzieren kann man das Loglevel entsprechend setzen:
Die Priorität ist wie folgt (höchste zu niedrigster):
```
NONE
WARN
INFO
DEBUG
VERBOSE
```

Wird das LogLevel also z.B. wie folgt auf Info gesetzt
```go
package main
import (
    log "intern.t3debug.de/git/sg-go-libs/AwesomeLog"
)

func main() {
    log.SetLogLevel(log.INFO)

    log.Println(log.VERBOSE, "Foobar Verbose")
    log.Println(log.DEBUG, "Foobar Debug")
    log.Println(log.INFO, "Foobar Info")
    log.Println(log.WARN, "Foobar Warning")
}
```
reduziert sich die Ausgabe auf:
![cmdline output](https://intern.t3debug.de/git/sg-go-libs/AwesomeLog/-/wikis/uploads/35e1f2ce6e8f688d7fc56b53d3e7e680/Auswahl_307.png)


##Notes
PrettyPrint kann nur Felder anzeigen die auch exportiert werden!
