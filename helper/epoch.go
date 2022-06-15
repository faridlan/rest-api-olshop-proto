package helper

import "time"

func EpochTime() int64 {
	now := time.Now()
	return now.UnixMilli()
}
