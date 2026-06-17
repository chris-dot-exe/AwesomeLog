package caller

import log "github.com/chris-dot-exe/AwesomeLog"

func Test() {
	log.Println(log.DEBUG, "Caller Test!")
}

func TestNew() {
	log.Println(log.DEBUG, "Caller Test!")
	log.Debug("Caller Test!")
}
