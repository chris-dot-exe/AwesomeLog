package test

import (
	"testing"

	log "github.com/chris-dot-exe/AwesomeLog"
	"github.com/chris-dot-exe/AwesomeLog/test/this/is/a/really/long/and/deep/nested/directory/to/test/the/caller"
)

func TestDeepNestedCaller(t *testing.T) {
	log.SetLogLevel(log.DEBUG)
	log.ShowColors(true)
	log.ShowColorsInLogs(true)
	log.ShowTimestamp(true)
	log.SetCallerMaxDepth(5)

	caller.Test()
}
