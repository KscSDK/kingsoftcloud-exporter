package util

import (
	"errors"
	"strings"
	"time"
)

const (
	TIME_TEMPLATE_UTC = "2006-01-02T15:04:05Z"

	ERROR_TIME = "time type should be 'yyyy-MM-ddTHH:mm:ssZ', like '2016-05-11T15:00:00Z'"
)

//ParseUTCToUnix
func ParseUTCToUnix(timeParse string) (int64, error) {

	timestamp, err := time.ParseInLocation(TIME_TEMPLATE_UTC, timeParse, time.Local)
	if err == nil {
		return timestamp.Unix(), nil
	}

	timeDay := strings.Split(timeParse, "T")
	if len(timeDay) != 2 {
		return 0, errors.New(ERROR_TIME)
	}

	patternDate := "2006-01-02"
	dayOverTime, err := time.ParseInLocation(patternDate, timeDay[0], time.Local)
	if err != nil {
		return 0, errors.New(ERROR_TIME)
	}

	hms := timeDay[1]
	if strings.HasSuffix(timeDay[1], "z") || strings.HasSuffix(timeDay[1], "Z") {
		hms = timeDay[1][:len(timeDay[1])-1]
	}

	dayInTime, err := dayInTimeParse(hms)
	if err != nil {
		return 0, errors.New(ERROR_TIME)
	}

	timeAll := dayOverTime.Add(dayInTime)
	return timeAll.Unix(), nil
}

func dayInTimeParse(timeParse string) (time.Duration, error) {
	timeParts := strings.Split(timeParse, ":")
	if len(timeParts) != 3 {
		return 0, errors.New(ERROR_TIME)
	}
	var timeDuring time.Duration
	timeDuringH, err := time.ParseDuration(timeParts[0] + "h")
	if err != nil {
		return 0, errors.New(ERROR_TIME)
	}
	timeDuringM, err := time.ParseDuration(timeParts[1] + "m")
	if err != nil {
		return 0, errors.New(ERROR_TIME)
	}
	timeDuringS, err := time.ParseDuration(timeParts[2] + "s")
	if err != nil {
		return 0, errors.New(ERROR_TIME)
	}
	timeDuring = timeDuringH + timeDuringM + timeDuringS
	return timeDuring, nil
}
