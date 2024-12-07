package timeconvert

import (
	"time"

	"github.com/RedEpochDB/RedEpochDB/lib/numconvert"
)

// (integer) UNIX sample timestamp in milliseconds. * can be used for an automatic timestamp from the system clock.
func ArgTimeParsing(arg []byte) (timestamp int64) {
	if string(arg) == "*" {
		timestamp = time.Now().UnixNano() / 1e6
	} else {
		timestamp = numconvert.BytesToInt64(string(arg))
	}
	return
}
