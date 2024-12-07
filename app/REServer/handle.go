package main

import (
	"math"
	"strings"
	"sync"

	"github.com/RedEpochDB/RedEpochDB/core/redhub"
	"github.com/RedEpochDB/RedEpochDB/core/redhub/pkg/resp"
	"github.com/RedEpochDB/RedEpochDB/lib/numconvert"
	"github.com/RedEpochDB/RedEpochDB/lib/timeconvert"
	tstorage "github.com/nakabonne/tstorage"
)

var mu sync.RWMutex

var items = make(map[string][]byte)

func (rts *REServer) NewRTSHandle() *redhub.RedHub {
	return redhub.NewRedHub(rts.onOpened, rts.onClose, rts.handle)
}

func (rts *REServer) onOpened(c *redhub.Conn) (out []byte, action redhub.Action) {
	return
}

func (rts *REServer) onClose(c *redhub.Conn, err error) (action redhub.Action) {
	return
}

func (rts *REServer) handle(cmd resp.Command, out []byte) ([]byte, redhub.Action) {
	var status redhub.Action
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		out = resp.AppendError(out, "ERR unknown command '"+string(cmd.Args[0])+"'")
	case "ts.add":
		//Append a sample to a time series.
		//TS.ADD key timestamp value [RETENTION retentionPeriod] [ENCODING [COMPRESSED|UNCOMPRESSED]] [CHUNK_SIZE size] [ON_DUPLICATE policy] [LABELS {label value}...]
		lens := len(cmd.Args)
		if lens < 4 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}
		timestamp := timeconvert.ArgTimeParsing(cmd.Args[2])
		rows := []tstorage.Row{{Metric: string(cmd.Args[1]), DataPoint: tstorage.DataPoint{
			Timestamp: timestamp, Value: numconvert.StringToFloat64(string(cmd.Args[3]))}}}
		for i := 4; i < lens; i++ {
			if strings.ToLower(string(cmd.Args[i])) == "labels" {
				rows[0].Labels = make([]tstorage.Label, 0)
				for j := i + 1; j+2 <= lens; j += 2 {
					rows[0].Labels = append(rows[0].Labels, tstorage.Label{
						Name:  string(cmd.Args[j]),
						Value: string(cmd.Args[j+1]),
					})
				}
				break
			}
		}
		if err := rts.MemoryDB.InsertRows(rows); err != nil {
			out = resp.AppendError(out, "ERR TS.ADD '"+err.Error()+"'")
		} else {
			out = resp.AppendInt(out, timestamp)
		}
	case "ts.madd":
		//Append new samples to one or more time series.
		//TS.MADD {key timestamp value}...
		lens := len(cmd.Args)
		if lens < 4 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}
		rows := make([]tstorage.Row, 0, 10)
		rows = append(rows, tstorage.Row{Metric: string(cmd.Args[1]), DataPoint: tstorage.DataPoint{
			Timestamp: timeconvert.ArgTimeParsing(cmd.Args[2]), Value: numconvert.StringToFloat64(string(cmd.Args[3]))}})
		for i := 4; i+3 <= lens; i += 3 {
			rows = append(rows, tstorage.Row{Metric: string(cmd.Args[i]), DataPoint: tstorage.DataPoint{
				Timestamp: timeconvert.ArgTimeParsing(cmd.Args[i+1]), Value: numconvert.StringToFloat64(string(cmd.Args[i+2]))}})
		}
		if err := rts.MemoryDB.InsertRows(rows); err != nil {
			out = resp.AppendError(out, "ERR TS.ADD '"+err.Error()+"'")
		} else {
			timestamp := make([]int64, 1)
			for _, v := range rows {
				timestamp = append(timestamp, v.Timestamp)
			}
			out = resp.AppendAny(out, timestamp)
		}
	case "ts.get":
		//Get the last sample.
		//TS.GET key
		if len(cmd.Args) != 2 {
			out = resp.AppendError(out, "ERR wrong number of arguments for '"+string(cmd.Args[0])+"' command")
			break
		}
		start := int64(0) // 你需要根据实际情况设置 start 的值
		end := math.MaxInt64

		sRow, err := rts.MemoryDB.Select(string(cmd.Args[1]), nil, start, int64(end))
		if err != nil {
			out = resp.AppendError(out, "ERR TS.RANGE '"+err.Error()+"'")
		}

		// Find the data point with the highest timestamp
		var latestPoint *tstorage.DataPoint
		for _, point := range sRow {
			if latestPoint == nil || point.Timestamp > latestPoint.Timestamp {
				latestPoint = point
			}
		}

		// Format the response

		dataMap := make(map[resp.SimpleInt64]float64)
		// for _, point := range sRow {
		// 	dataMap[resp.SimpleInt64(point.Timestamp)] = point.Value
		// }

		if latestPoint != nil {
			// Return an array of timestamp and value
			dataMap[resp.SimpleInt64(latestPoint.Timestamp)] = latestPoint.Value
		}
		out = resp.AppendAny(out, dataMap)

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
		dataMaps := make([]map[resp.SimpleInt64]float64, 0)
		for _, v := range sRows {
			dataMaps = append(dataMaps, map[resp.SimpleInt64]float64{resp.SimpleInt64(v.Timestamp): v.Value})
		}
		out = resp.AppendAny(out, dataMaps)
	case "ts.alter":
		//Update the retention, chunk size, duplicate policy, and labels of an existing time series.
		//TS.ALTER key [RETENTION retentionPeriod] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS [{label value}...]]
		out = resp.AppendString(out, "OK")
	case "ts.create":
		//Create a new time series.
		//TS.CREATE key [RETENTION retentionPeriod] [ENCODING [UNCOMPRESSED|COMPRESSED]] [CHUNK_SIZE size] [DUPLICATE_POLICY policy] [LABELS {label value}...]
		out = resp.AppendString(out, "OK")
	case "ts.createrule":
		//Create a compaction rule.
		//TS.CREATERULE sourceKey destKey AGGREGATION aggregator bucketDuration
		out = resp.AppendString(out, "OK")
	case "ts.decrby":
		//Decrease the value of the sample with the maximal existing timestamp, or create a new sample with a value equal to the value of the sample with the maximal existing timestamp with a given decrement.
		//TS.DECRBY key value [TIMESTAMP timestamp] [RETENTION retentionPeriod] [UNCOMPRESSED] [CHUNK_SIZE size] [LABELS {label value}...]
		out = resp.AppendString(out, "OK")
	case "ts.del":
		//Delete all samples between two timestamps for a given time series.
		//The given timestamp interval is closed (inclusive), meaning samples which timestamp eqauls the fromTimestamp or toTimestamp will also be deleted.
		//TS.DEL key fromTimestamp toTimestamp
		out = resp.AppendString(out, "OK")
	case "ts.deleterule":
		//Delete a compaction rule.
		//TS.DELETERULE sourceKey destKey
		out = resp.AppendString(out, "OK")
	case "ts.incrby":
		//Increase the value of the sample with the maximal existing timestamp, or create a new sample with a value equal to the value of the sample with the maximal existing timestamp with a given increment.
		//TS.INCRBY key value [TIMESTAMP timestamp] [RETENTION retentionPeriod] [UNCOMPRESSED] [CHUNK_SIZE size] [LABELS {label value}...]
		out = resp.AppendString(out, "OK")
	case "ts.info":
		//Returns information and statistics for a time series.
		//TS.INFO key [DEBUG]
		out = resp.AppendString(out, "OK")
	case "ts.mget":
		//Get the last samples matching a specific filter.
		//TS.MGET [WITHLABELS | SELECTED_LABELS label...] FILTER filter...
		out = resp.AppendString(out, "OK")
	case "ts.mrange":
		//Query a range across multiple time series by filters in forward direction.
		//TS.MRANGE fromTimestamp toTimestamp [FILTER_BY_TS TS...] [FILTER_BY_VALUE min max] [WITHLABELS | SELECTED_LABELS label...] [COUNT count] [ALIGN value] [AGGREGATION aggregator bucketDuration] FILTER filter.. [GROUPBY label REDUCE reducer]
		out = resp.AppendString(out, "OK")
	case "ts.mrevrange":
		//Query a range across multiple time series by filters in reverse direction.
		//TS.MREVRANGE fromTimestamp toTimestamp [FILTER_BY_TS TS...] [FILTER_BY_VALUE min max] [WITHLABELS | SELECTED_LABELS label...] [COUNT count] [ALIGN value] [AGGREGATION aggregator bucketDuration] FILTER filter.. [GROUPBY label REDUCE reducer]
		out = resp.AppendString(out, "OK")
	case "ts.queryindex":
		//Get all time series keys matching a filter list.
		//TS.QUERYINDEX filter...
		out = resp.AppendString(out, "OK")
	case "ts.revrange":
		//Query a range in reverse direction.
		//TS.REVRANGE key fromTimestamp toTimestamp [FILTER_BY_TS TS...] [FILTER_BY_VALUE min max] [COUNT count] [ALIGN value] [AGGREGATION aggregator bucketDuration]
		out = resp.AppendString(out, "OK")
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
