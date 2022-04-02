/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:37:18
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 13:34:00
 * @FilePath: /RedTimeDB/app/rtserver/main.go
 */

package main

import (
	"log"
)

const configFileDefaultName = "rtserver.yaml"

//BuildDate: Binary file compilation time
//BuildVersion: Binary compiled GIT version
var (
	BuildDate    string
	BuildVersion string
)

func main() {
	// parse configure
	var configFileURI string
	parseCmdArgs(&configFileURI) //Parse runtime parameters.

	// init new RTServer object
	rtserver, err := GetNewRTServer(configFileURI)

	if err != nil {
		log.Panic(err)
	}

	if rtserver.RH != nil {
		rtserver.ListendAndServe()
	}

	return
}
