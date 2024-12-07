package main

import (
	"sync"

	confer "github.com/RedEpochDB/RedEpochDB/app/REServer/REServerconf"
	"github.com/RedEpochDB/RedEpochDB/core/redhub"
	"github.com/RedEpochDB/RedEpochDB/lib/gnet"
	tstorage "github.com/nakabonne/tstorage"
)

// REServer is the core structure of RedEpochDB
type REServer struct {
	Confer   confer.Confer
	Mutex    sync.RWMutex
	RH       *redhub.RedHub
	MemoryDB tstorage.Storage
}

// GetNewREServer Used to create the RedEpochDB core structure
func GetNewREServer(conFileURI string) (REServer REServer, err error) {
	//get confer object
	REServer.Confer, err = confer.GetNewConfer(conFileURI)
	if err != nil {
		return
	}

	REServer.MemoryDB, err = tstorage.NewStorage(
		tstorage.WithTimestampPrecision(tstorage.Milliseconds),
	)
	if err != nil {
		return
	}
	//Focus adjustment section
	REServer.RH = REServer.NewRTSHandle()

	return
}

func (rts *REServer) ListendAndServe() error {
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
