# AwesomeLog

AwesomeLog ist eine inplace Erweiterung des Standard log packages und bietet die Möglichkeit verschiedene Log Level zu definieren.

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

Implementierte Funktionen sind `log.Println()`, `log.Print()`, `log.Printf()`, `log.PrettyPrint()` sowie
`log.SetLogLevel()`, `log.SetDefaultLevel()` und `log.ShowColorsInLogs()`

Die Print-Funktionen sind in denn Beispielen unten zu sehen.

Folgend aber kurz die Konfigurationsfunktionen erläutert:
#### `log.SetLogLevel(logLevel)`
`log.SetLogLevel()` definiert bis zu welchem Log Level Log-Nachrichten ausgegeben werden sollen.

#### `log.SetLogLevelByString(string)`
`log.SetLogLevelByString()` ist identisch mit `SetLogLevel()` nur das hier das Level als String (z.B. aus einem Config File) angegeben wird.

Default: log.NONE

Siehe Examples unten.

#### `log.SetDefaultLevel(logLevel)`
`log.SetDefaultLevel()` definiert welches Log-Level benutzt wird wenn kein Log-Level angegeben wurde.

Default: log.INFO

#### `log.ShowColorsInLogs(bool)`
`log.ShowColorsInLogs()` definiert ob die Farbigen Log-Level Labels auch mit ausgegeben werden sollen wenn das Log nicht im Terminal ausgegeben wird sondern in ein File umgeleitet wird.
Dies funktioniert auch wenn die Log-Nachrichten im Terminal ausgegeben UND in ein File gespeichert werden (Docker Logs z.B.)

Default: false

Nicht ausreichend getestet, kann also zu Fehlern führen.

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

![cmdline output](https://user-images.githubusercontent.com/49272981/80649110-b3297e80-8a71-11ea-9779-d359da872d75.png)



### Examples with Loglevel
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

![cmdline output](https://user-images.githubusercontent.com/49272981/80649108-b290e800-8a71-11ea-8463-595e9de9f171.png)

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
![cmdline output](https://user-images.githubusercontent.com/49272981/80649104-b1f85180-8a71-11ea-9c72-d98ed825a6b4.png)


## Notes
PrettyPrint kann nur Felder anzeigen die auch exportiert werden!
