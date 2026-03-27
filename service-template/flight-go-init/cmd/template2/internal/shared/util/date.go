package util

import (
	"fmt"
	"strings"
	"time"
)

const (
	DateFormat                    = "2006-01-02"
	DateTimeFormatWithSecond      = "2006-01-02 15:04:05"
	DateTimeFormat                = "2006-01-02 15:04"
	SegmentDateTimeFormat         = DateTimeFormatRFC3339Zoneless
	DateTimeFormatRFC3339Zoneless = "2006-01-02T15:04:05"
	DateTimeFormatWithoutSecond   = "2006-01-02 15:04"
	TimeFormat                    = "15:04"
)

func ParseDuration(s string, value time.Duration) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return value
	}

	return d
}

func ParseTime(format string, s string) *time.Time {
	t, err := time.Parse(format, s)
	if err != nil {
		return nil
	}
	return &t
}

func FormatTime(format string, t time.Time) string {
	s := t.Format(format)
	return s
}

// SplitDateTime 2024-11-11T12:10:00 -> 2024-11-11, 12:10
func SplitDateTime(dateTimeStr string) (string, string) {
	// Split the date and time using "T" as the delimiter
	parts := strings.Split(dateTimeStr, "T")
	if len(parts) != 2 {
		return "", ""
	}

	// parts[0] is the date, parts[1] is the time with seconds
	date := parts[0]
	if date == "" {
		return "", ""
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) < 2 {
		return "", ""
	}

	// Construct the time string without seconds
	time := fmt.Sprintf("%s:%s", timeParts[0], timeParts[1])

	return date, time
}

func ConvertAndFormatTimezone(dateStr, sourceFormat string, targetFormat string, timezoneLoc string) (string, error) {
	// Parse the formatted string back to time in UTC
	parsedTime, err := time.ParseInLocation(sourceFormat, dateStr, time.UTC)
	if err != nil {
		return "", err
	}
	// Load the specified timezone location
	location, err := time.LoadLocation(timezoneLoc)
	if err != nil {
		return "", fmt.Errorf("failed to load timezone location: %w", err)
	}

	// Convert the parsed time to the specified timezone
	localTime := parsedTime.In(location)

	// Format the parsed time to the desired format
	return localTime.Format(targetFormat), nil
}
