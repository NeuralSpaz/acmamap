package main

import (
	"log"
	"log/syslog"
)

func intt() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "ACMAMAP")
	if e == nil {
		log.SetOutput(logwriter)
	}
}

func main() {
	log.Println("Starting ACMAMAP")
}
