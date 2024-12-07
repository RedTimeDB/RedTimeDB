package redhub

import (
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/RedEpochDB/RedEpochDB/lib/gnet"

	"github.com/RedEpochDB/RedEpochDB/core/redhub/pkg/resp"
)

type Action int

const (
	// None indicates that no action should occur following an event.
	None Action = iota

	// Close closes the connection.
	Close

	// Shutdown shutdowns the server.
	Shutdown
)

type Conn struct {
	gnet.Conn
}

type Options struct {
	// Multicore indicates whether the server will be effectively created with multi-cores, if so,
	// then you must take care with synchronizing memory between all event callbacks, otherwise,
	// it will run the server with single thread. The number of threads in the server will be automatically
	// assigned to the value of logical CPUs usable by the current process.
	Multicore bool

	// LockOSThread is used to determine whether each I/O event-loop is associated to an OS thread, it is useful when you
	// need some kind of mechanisms like thread local storage, or invoke certain C libraries (such as graphics lib: GLib)
	// that require thread-level manipulation via cgo, or want all I/O event-loops to actually run in parallel for a
	// potential higher performance.
	LockOSThread bool

	// ReadBufferCap is the maximum number of bytes that can be read from the client when the readable event comes.
	// The default value is 64KB, it can be reduced to avoid starving subsequent client connections.
	//
	// Note that ReadBufferCap will be always converted to the least power of two integer value greater than
	// or equal to its real amount.
	ReadBufferCap int

	// LB represents the load-balancing algorithm used when assigning new connections.
	LB gnet.LoadBalancing

	// NumEventLoop is set up to start the given number of event-loop goroutine.
	// Note: Setting up NumEventLoop will override Multicore.
	NumEventLoop int

	// ReusePort indicates whether to set up the SO_REUSEPORT socket option.
	ReusePort bool

	// Ticker indicates whether the ticker has been set up.
	Ticker bool

	// TCPKeepAlive sets up a duration for (SO_KEEPALIVE) socket option.
	TCPKeepAlive time.Duration

	// TCPNoDelay controls whether the operating system should delay
	// packet transmission in hopes of sending fewer packets (Nagle's algorithm).
	//
	// The default is true (no delay), meaning that data is sent
	// as soon as possible after a Write.
	TCPNoDelay gnet.TCPSocketOpt

	// SocketRecvBuffer sets the maximum socket receive buffer in bytes.
	SocketRecvBuffer int

	// SocketSendBuffer sets the maximum socket send buffer in bytes.
	SocketSendBuffer int

	// ICodec encodes and decodes TCP stream.
	Codec gnet.ICodec
}

func NewRedHub(
	onOpened func(c *Conn) (out []byte, action Action),
	onClosed func(c *Conn, err error) (action Action),
	handler func(cmd resp.Command, out []byte) ([]byte, Action),
) *RedHub {
	return &RedHub{
		redHubBufMap: make(map[gnet.Conn]*connBuffer),
		connSync:     sync.RWMutex{},
		onOpened:     onOpened,
		onClosed:     onClosed,
		handler:      handler,
	}
}

type RedHub struct {
	*gnet.EventServer
	onOpened     func(c *Conn) (out []byte, action Action)
	onClosed     func(c *Conn, err error) (action Action)
	handler      func(cmd resp.Command, out []byte) ([]byte, Action)
	redHubBufMap map[gnet.Conn]*connBuffer
	connSync     sync.RWMutex
}

type connBuffer struct {
	buf     bytes.Buffer
	command []resp.Command
}

func (rs *RedHub) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	rs.connSync.Lock()
	defer rs.connSync.Unlock()
	rs.redHubBufMap[c] = new(connBuffer)
	rs.onOpened(&Conn{Conn: c})
	return
}

func (rs *RedHub) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	rs.connSync.Lock()
	defer rs.connSync.Unlock()
	delete(rs.redHubBufMap, c)
	rs.onClosed(&Conn{Conn: c}, err)
	return
}

func (rs *RedHub) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	//decodebuff := credis.NewDecoder(c.Respbuffer())
	//for {
	//	command, err := decodebuff.Decode()
	//	if err != nil {
	//		//if err.Error() == io.EOF.Error() {
	//		//	action = gnet.Close
	//		//}
	//		return
	//	}
	//
	//	if command.Type != credis.TypeArray {
	//		return
	//	}
	//
	//	commandsCount := len(command.Array)
	//
	//	if commandsCount < 1 {
	//		return
	//	}
	//	if command.Array[0].Type != credis.TypeBulkBytes {
	//		return
	//	}
	//
	//	commandArgs := make([]interface{}, commandsCount)
	//	for i := 0; i < commandsCount; i++ {
	//		commandArgs[i] = command.Array[i].Value
	//	}
	//	var status Action
	//	out, status = rs.handler1(commandArgs, out)
	//	switch status {
	//	case Close:
	//		action = gnet.Close
	//	}
	//}
	//return

	r := resp.NewReader(c.Respbuffer())
	var status Action
	for {
		cmd, err, last := r.ReadCommand()
		if err != nil {
			if len(last) != 0 {
				c.BuffLock()
				c.Respbuffer().Write(last)
				c.BuffUnLock()
				break
			}
			if err == io.EOF {
				//action = gnet.Close
				break
			}
		}
		out, status = rs.handler(cmd, out)
		switch status {
		case Close:
			action = gnet.Close
		}
	}
	return
}

func ListendAndServe(addr string, options Options, rh *RedHub) error {
	serveOptions := gnet.Options{
		Multicore:        options.Multicore,
		LockOSThread:     options.LockOSThread,
		ReadBufferCap:    options.ReadBufferCap,
		LB:               options.LB,
		NumEventLoop:     options.NumEventLoop,
		ReusePort:        options.ReusePort,
		Ticker:           options.Ticker,
		TCPKeepAlive:     options.TCPKeepAlive,
		TCPNoDelay:       options.TCPNoDelay,
		SocketRecvBuffer: options.SocketRecvBuffer,
		SocketSendBuffer: options.SocketSendBuffer,
		Codec:            options.Codec,
	}

	return gnet.Serve(rh, addr, gnet.WithOptions(serveOptions))
}
