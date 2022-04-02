/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:42:01
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 11:43:58
 * @FilePath: /RedTimeDB/app/rtserver/common.go
 */
package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	showBanner()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
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
