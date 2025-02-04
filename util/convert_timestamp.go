package util

import (
	"fmt"
	"time"
)

func ConvertTimeStamp(timestamp string) (int64, error) {
	tm, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		fmt.Errorf("unable to convert the timestamp: ", err)
		return 0, err
	}

	seconds := tm.Unix()
	return seconds, nil
}
