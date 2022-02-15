/*
AwesomeLog is a fully compatible drop-in replacement for the standard library logger with some awesome features.

AwesomeLog let you define log levels for each logged messages as well as a PrettyPrint function to print out objects in a pretty readable format.
It also adds the option to show details of the caller like file path, function name and line number.

Since version v0.11.0 AwesomeLog also provides the functionality to add custom handlers for each LogLevel.

The simplest way to use AwesomeLog is to just replace the standard library log with AwesomeLog:

  package main

  import  log "github.com/chris-dot-exe/AwesomeLog"

  func main() {
    log.Println("Hello World!")
  }

*/
package log
