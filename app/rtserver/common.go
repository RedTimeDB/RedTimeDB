/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:42:01
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 12:16:13
 * @FilePath: /RedTimeDB/app/rtserver/common.go
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	showBanner()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func parseCmdArgs(configURI *string) {

	dir, err := os.Getwd()
	if err != nil {
		dir = configFileDefaultName
	} else {
		dir += "/" + configFileDefaultName
	}

	//Set the parsing parameter -c and set the default configuration file path to os.Getwd()/CONFIG_FILE_DEFAULT_NAME
	flag.StringVar(configURI, "c", dir, "Config file path")

	flag.Parse()
}

func showBanner() {
	bannerData := `         __                                
   _____/ /_________  ______   _____  _____
  / ___/ __/ ___/ _ \/ ___/ | / / _ \/ ___/
 / /  / /_(__  )  __/ /   | |/ /  __/ /    
/_/   \__/____/\___/_/    |___/\___/_/     
                                           
`
	fmt.Println(bannerData)
	fmt.Println("Build Version: ", BuildVersion, "  Date: ", BuildDate)
	time.Sleep(time.Millisecond * 100)
}
