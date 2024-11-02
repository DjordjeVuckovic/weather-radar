package util

import (
	"time"
)

func UnixToUTC(unix int64) string {
	return time.Unix(unix, 0).UTC().Format(time.RFC3339)
}

func UnixToLocal(unix int64, tzOffset int16) string {
	t := time.Unix(unix, 0).UTC()

	localTime := t.Add(time.Duration(tzOffset) * time.Second)

	return localTime.Format("2006-01-02 15:04")
}
