package restapi

import "time"

func timeToEpochMs(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / 1000000
}

func epochMsToTime(epochMs int64) time.Time {
	return time.Unix(epochMs/1000, epochMs%1000*1000000).UTC()
}
