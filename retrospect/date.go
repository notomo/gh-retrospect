package retrospect

import "time"

const yyyymmdd = "2006-01-02"

func ParseFrom(from string) (time.Time, error) {
	duration, err := time.ParseDuration(from)
	if err == nil {
		return time.Now().Add(duration), nil
	}

	if from == "" {
		now := time.Now().AddDate(0, 0, -7)
		from = now.Format(yyyymmdd)
	}
	return time.ParseInLocation(yyyymmdd, from, time.Local)
}

func ParseTo(to string) (time.Time, error) {
	if to == "" {
		return time.Time{}, nil
	}

	duration, err := time.ParseDuration(to)
	if err == nil {
		return time.Now().Add(duration), nil
	}

	return time.ParseInLocation(yyyymmdd, to, time.Local)
}
