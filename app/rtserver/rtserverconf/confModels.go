/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:56:55
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 13:18:47
 * @FilePath: /RedTimeDB/app/rtserver/rtserverconf/confModels.go
 */

package confer

import (
	"sync"
)

type Confer struct {
	Mutex sync.RWMutex
	Opts  RTServerConfS
}

type RTServerConfS struct {
	NetConf   NetConfS   `yaml:"net"`
	ApiConf   ApiConfS   `yaml:"api"`
	DebugConf DebugConfS `yaml:"debug"`
	CacheConf CacheS     `yaml:"cache"` //Cache related configuration
}

//CacheS is cache configuration
type CacheS struct {
	Enable            bool `yaml:"enable"`
	MaxItemsSize      int  `yaml:"max_items_size"`     //Maximum number of items in cache
	DefaultExpiration int  `yaml:"default_expiration"` //The default expiration time of ache kv (unit: milliseconds)
	CleanupInterval   int  `yaml:"cleanup_interval"`   //Cache expired kv cleaning cycle (unit: seconds)
}

type DebugConfS struct {
	Enable   bool   `yaml:"enable"`
	PprofUri string `yaml:"pprof_uri"` //Performance Monitoring - Listening Address
}

type NetConfS struct {
	ListenUri string `yaml:"listen_uri"` //network listening address
}

type ApiConfS struct {
	HttpListenAddress string `yaml:"http_listen_address"` //api listening address
}
