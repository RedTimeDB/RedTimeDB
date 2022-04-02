/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:53:31
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 13:51:55
 * @FilePath: /RedTimeDB/app/rtserver/rtserver.go
 */

package main

import (
	"strings"
	"sync"

	confer "github.com/RedTimeDB/RedTimeDB/app/rtserver/rtserverconf"
	"github.com/RedTimeDB/RedTimeDB/core/redhub"
	"github.com/RedTimeDB/RedTimeDB/core/redhub/pkg/resp"
)

//RTServer is the core structure of RedTimeServer
type RTServer struct {
	Confer confer.Confer
	Mutex  sync.RWMutex
	RH     *redhub.RedHub
}

//GetNewRTServer Used to create the RedTimeServer core structure
func GetNewRTServer(conFileURI string) (rtserver RTServer, err error) {
	//get confer object
	rtserver.Confer, err = confer.GetNewConfer(conFileURI)
	if err != nil {
		return
	}

	//Focus adjustment section
	var mu sync.RWMutex
	var items = make(map[string][]byte)
	rtserver.RH = redhub.NewRedHub(
		func(c *redhub.Conn) (out []byte, action redhub.Action) {
			return
		},
		func(c *redhub.Conn, err error) (action redhub.Action) {
			return
		},
		func(cmd resp.Command, out []byte) ([]byte, redhub.Action) {
			var status redhub.Action
			switch strings.ToLower(string(cmd.Args[0])) {
			default:
				out = resp.AppendError(out, "ERR unknown command '"+string(cmd.Args[0])+"'")
			case "ping":
				out = resp.AppendString(out, "PONG")
			case "quit":
				out = resp.AppendString(out, "OK")
				status = redhub.Close
			case "set":
				if len(cmd.Args) != 3 {
					out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
					break
				}
				mu.Lock()
				items[string(cmd.Args[1])] = cmd.Args[2]
				mu.Unlock()
				out = resp.AppendString(out, "OK")
			case "get":
				if len(cmd.Args) != 2 {
					out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
					break
				}
				mu.RLock()
				val, ok := items[string(cmd.Args[1])]
				mu.RUnlock()
				if !ok {
					out = resp.AppendNull(out)
				} else {
					out = resp.AppendBulk(out, val)
				}
			case "del":
				if len(cmd.Args) != 2 {
					out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
					break
				}
				mu.Lock()
				_, ok := items[string(cmd.Args[1])]
				delete(items, string(cmd.Args[1]))
				mu.Unlock()
				if !ok {
					out = resp.AppendInt(out, 0)
				} else {
					out = resp.AppendInt(out, 1)
				}
			case "config":
				// This simple (blank) response is only here to allow for the
				// redis-benchmark command to work with this example.
				out = resp.AppendArray(out, 2)
				out = resp.AppendBulk(out, cmd.Args[2])
				out = resp.AppendBulkString(out, "")
			}
			return out, status
		},
	)

	return
}

func (rts *RTServer) ListendAndServe(addr string, options redhub.Options) {

}
