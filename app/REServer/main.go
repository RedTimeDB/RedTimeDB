package main

import (
	"log"
)

const configFileDefaultName = "REServer.yaml"

// BuildDate: Binary file compilation time
// BuildVersion: Binary compiled GIT version
var (
	BuildDate    string
	BuildVersion string
)

func main() {
	// parse configure
	var configFileURI string
	parseCmdArgs(&configFileURI) //Parse runtime parameters.

	// init new REServer object
	REServer, err := GetNewREServer(configFileURI)

	if err != nil {
		log.Panic(err)
	}

	if REServer.RH != nil {
		if err = REServer.ListendAndServe(); err != nil {
			log.Panic(err)
		}
	}

	return
}
