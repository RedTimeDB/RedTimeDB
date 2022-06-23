/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:53:31
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 13:51:55
 * @FilePath: /RedTimeDB/app/rtserver/rtserver.go
 */

package main

import (
	"sync"

	"github.com/RedTimeDB/RedTimeDB/lib/gnet"

	confer "github.com/RedTimeDB/RedTimeDB/app/rtserver/rtserverconf"
	"github.com/RedTimeDB/RedTimeDB/core/redhub"
	tstorage "github.com/RedTimeDB/RedTimeDB/core/storage"
)

//RTServer is the core structure of RedTimeServer
type RTServer struct {
	Confer   confer.Confer
	Mutex    sync.RWMutex
	RH       *redhub.RedHub
	MemoryDB tstorage.Storage
}

//GetNewRTServer Used to create the RedTimeServer core structure
func GetNewRTServer(conFileURI string) (rtserver RTServer, err error) {
	//get confer object
	rtserver.Confer, err = confer.GetNewConfer(conFileURI)
	if err != nil {
		return
	}

	rtserver.MemoryDB, err = tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Milliseconds),
	)
	if err != nil {
		return
	}
	//Focus adjustment section
	rtserver.RH = rtserver.NewRTSHandle()

	return
}

func (rts *RTServer) ListendAndServe() error {
	serveOptions := gnet.Options{
		Multicore:        rts.Confer.Opts.NetConf.Muticore,
		LockOSThread:     rts.Confer.Opts.NetConf.LockOSThread,
		ReadBufferCap:    rts.Confer.Opts.NetConf.ReadBufferCap,
		LB:               gnet.LoadBalancing(rts.Confer.Opts.NetConf.LB),
		NumEventLoop:     rts.Confer.Opts.NetConf.NumEventLoop,
		ReusePort:        rts.Confer.Opts.NetConf.ReusePort,
		Ticker:           rts.Confer.Opts.NetConf.Ticker,
		TCPKeepAlive:     rts.Confer.Opts.NetConf.TCPKeepAlive,
		TCPNoDelay:       gnet.TCPSocketOpt(rts.Confer.Opts.NetConf.TCPNoDelay),
		SocketRecvBuffer: rts.Confer.Opts.NetConf.SocketRecvBuffer,
		SocketSendBuffer: rts.Confer.Opts.NetConf.SocketSendBuffer,
	}

	return gnet.Serve(rts.RH, rts.Confer.Opts.NetConf.ListenUri, gnet.WithOptions(serveOptions))
}
