package repository

import (
	"encoding/base64"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999999Z00:00" // reduce precision from RFC3339Nano as date format
)

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}
	timeString := string(byt)
	if len(encodedTime) == 0 {
		timeString = time.Now().Format(timeFormat)
	}
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
