package helper

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertDateTimeStringToOtherFormat(dateTimeString, dateTimeFormat, targetFormat string) (string, error) {
	// Parse the dateTimeString according to the provided dateTimeFormat
	parsedTime, err := time.Parse(dateTimeFormat, dateTimeString)
	if err != nil {
		return "", err
	}

	// Format the parsed time to the target format
	formattedTime := parsedTime.Format(targetFormat)

	return formattedTime, nil
}

// FormatTime receive HHMM to HH:MM format
func FormatTime(timeStr string) string {
	// Ensure the string is the correct length to prevent a panic
	if len(timeStr) != 4 {
		// try again use PadTime, sometimes, {{PROVIDER}} return 3 characters time
		return FormatTime(PadTime(timeStr))
	}
	// Slice the string and insert the colon
	return timeStr[0:2] + ":" + timeStr[2:4]
}

// PadTime adds a leading zero to a time string if it's only 3 characters long.
// It ensures the output is always in "HHMM" format (4 characters).
func PadTime(timeStr string) string {
	// If the string is already the correct length, do nothing.
	if len(timeStr) >= 4 {
		return timeStr
	}

	// Convert the string to an integer to handle the formatting.
	num, err := strconv.Atoi(timeStr)
	if err != nil {
		// If it's not a valid number, return the original string to avoid crashing.
		return timeStr
	}

	// Sprintf with "%04d" formats the integer:
	// d: format as a decimal integer
	// 4: pad to a width of 4 characters
	// 0: pad with leading zeros instead of spaces
	return fmt.Sprintf("%04d", num)
}

// input  20250717 13:53:11
// output "2025-07-17 13:53:11"
func FormatDateTimeStr(input string) string {
	parsedTime, err := time.Parse("20060102 15:04:05", input)
	if err != nil {
		// Optional: log or handle error properly
		return input
	}
	return parsedTime.Format("2006-01-02 15:04:05")
}

func ConvertToPlanTime(dateStr, timeStr string) (string, error) {
	parsedDate, err := time.Parse("20060102", dateStr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", parsedDate.Format("2006-01-02"), timeStr), nil
}
