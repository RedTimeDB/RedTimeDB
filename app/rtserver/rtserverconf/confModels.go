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
	"time"
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
	ListenUri        string        `yaml:"listen_uri"` //network listening address
	Muticore         bool          `yaml:"muticore"`
	NumEventLoop     int           `yaml:"num_event_loop"`
	LB               int           `yaml:"lb"`
	ReuseAddr        int           `yaml:"reuse_addr"`
	ReusePort        bool          `yaml:"reuse_port"`
	ReadBufferCap    int           `yaml:"read_buffer_cap"`
	LockOSThread     bool          `yaml:"lock_os_thread"`
	Ticker           bool          `yaml:"ticker"`
	TCPKeepAlive     time.Duration `yaml:"tcp_keep_alive"`
	TCPNoDelay       int           `yaml:"tcp_no_delay"`
	SocketRecvBuffer int           `yaml:"socket_recv_buffer"` // SocketRecvBuffer sets the maximum socket receive buffer in bytes.
	SocketSendBuffer int           `yaml:"socket_send_buffer"` // SocketSendBuffer sets the maximum socket send buffer in bytes.
}

type ApiConfS struct {
	HttpListenAddress string `yaml:"http_listen_address"` //api listening address
}

type MemoryDb struct {
	DataPath           string `yaml:"data_path"`
	PartitionDuration  string `yaml:"partition_duration"`
	Retention          string `yaml:"retention"`
	TimestampPrecision string `yaml:"timestamp_precision"`
	WriteTimeout       string `yaml:"write_timeout"`
	WALBufferedSize    int    `yaml:"wal_buffered_size"`
}
