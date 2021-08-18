package pkg

import (
	"time"
)

// ParseTime converts string time to time.Time struct.
func ParseTime(v string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
