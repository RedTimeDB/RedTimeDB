package main

import (
	"strings"
	"sync"
	"time"

	"github.com/RedTimeDB/RedTimeDB/core/redhub"
	"github.com/RedTimeDB/RedTimeDB/core/redhub/pkg/resp"
	tstorage "github.com/RedTimeDB/RedTimeDB/core/storage"
	"github.com/RedTimeDB/RedTimeDB/lib/numconvert"
)

var mu sync.RWMutex

var items = make(map[string][]byte)

func (rts *RTServer) NewRTSHandle() *redhub.RedHub {
	return redhub.NewRedHub(rts.onOpened, rts.onClose, rts.handle)
}

func (rts *RTServer) onOpened(c *redhub.Conn) (out []byte, action redhub.Action) {
	return
}

func (rts *RTServer) onClose(c *redhub.Conn, err error) (action redhub.Action) {
	return
}

func (rts *RTServer) handle(cmd resp.Command, out []byte) ([]byte, redhub.Action) {
	var status redhub.Action
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		out = resp.AppendError(out, "ERR unknown command '"+string(cmd.Args[0])+"'")
	case "ts.add":
		//Append a sample to a time series.
		//TS.ADD key timestamp value [RETENTION retentionPeriod] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [ON_DUPLICATE policy] [LABELS {label value}...]
		if len(cmd.Args) != 4 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}
		var timestamp int64
		if string(cmd.Args[2]) == "*" {
			timestamp = time.Now().UnixNano() / 1e6
		} else {
			timestamp = numconvert.BytesToInt64(string(cmd.Args[2]))
		}
		if err := rts.MemoryDB.InsertRows([]tstorage.Row{{Metric: string(cmd.Args[1]), DataPoint: tstorage.DataPoint{
			Timestamp: timestamp, Value: numconvert.StringToFloat64(string(cmd.Args[3]))}}}); err != nil {
			out = resp.AppendError(out, "ERR TS.ADD '"+err.Error()+"'")
		} else {
			out = resp.AppendInt(out, timestamp)
		}
	case "ts.get":
		//Get the last sample.
		//TS.GET key
		if len(cmd.Args) != 2 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}

	case "ts.range":
		//Get the last sample.
		//TS.GET key
		if len(cmd.Args) != 4 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}
		sRows, err := rts.MemoryDB.Select(string(cmd.Args[1]), nil, numconvert.BytesToInt64(string(cmd.Args[2])), numconvert.BytesToInt64(string(cmd.Args[3])))
		if err != nil {
			out = resp.AppendError(out, "ERR TS.RANGE '"+err.Error()+"'")
		}
		name := make([]map[resp.SimpleInt64]float64, 0)
		for _, v := range sRows {
			name = append(name, map[resp.SimpleInt64]float64{resp.SimpleInt64(v.Timestamp): v.Value})
		}
		out = resp.AppendAny(out, name)
	case "ts.alter":
		//Update the retention, chunk size, duplicate policy, and labels of an existing time series.
		//TS.ALTER key [RETENTION retentionPeriod] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS [{label value}...]]
	case "TS.CREATE":
		//Create a new time series.
		//TS.CREATE key [RETENTION retentionPeriod] [ENCODING [UNCOMPRESSED|COMPRESSED]] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS {label value}...]
	case "TS.CREATERULE":
		//Create a compaction rule.
		//TS.CREATERULE sourceKey destKey AGGREGATION aggregator bucketDuration
	case "TS.DECRBY":
		//Decrease the value of the sample with the maximal existing timestamp, or create a new sample with a value equal to the value of the sample with the maximal existing timestamp with a given decrement.
		//TS.DECRBY key value [TIMESTAMP timestamp] [RETENTION retentionPeriod] [UNCOMPRESSED] [CHUNK_SIZE size] [LABELS {label value}...]
	case "TS.MADD":
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
}
