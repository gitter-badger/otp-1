package logging

import (
	"log"
)

// Debug allows easy printing of debug logs
func Debug(msg string, debug bool) {
	if debug {
		log.Println(msg)
	} else {
		return
	}
}
