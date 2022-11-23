package tk

import "time"

const (
	TimeLayoutFormatMonth  = "2006-01"
	TimeLayoutFormatDay    = "2006-01-02"
	TimeLayoutFormatHour   = "2006-01-02 15"
	TimeLayoutFormatMinute = "2006-01-02 15:04"
	TimeLayoutFormatSecond = "2006-01-02 15:04:05"
)

var _timeLayoutFormats = []string{
	time.RFC3339,
	TimeLayoutFormatSecond,
	TimeLayoutFormatMinute,
	TimeLayoutFormatHour,
	TimeLayoutFormatDay,
	TimeLayoutFormatMonth,
}

func ParseTime(s string) (t time.Time, err error) {
	for _, format := range _timeLayoutFormats {
		t, err = time.ParseInLocation(format, s, time.Local)
		if err == nil {
			return
		}
	}
	return
}

func ParseTimePointer(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := ParseTime(s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
