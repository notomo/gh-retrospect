package retrospect

import "time"

const yyyymmdd = "2006-01-02"

func ParseDate(s string) (time.Time, error) {
	if s == "" {
		now := time.Now().AddDate(0, 0, -7)
		s = now.Format(yyyymmdd)
	}
	return time.ParseInLocation(yyyymmdd, s, time.Local)
}
