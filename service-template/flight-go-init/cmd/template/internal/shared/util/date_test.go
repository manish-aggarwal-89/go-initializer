package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDuration(t *testing.T) {
	// Test valid duration string
	durationStr := "1h30m"
	expectedDuration := 1*time.Hour + 30*time.Minute
	result := ParseDuration(durationStr, 0)
	if result != expectedDuration {
		t.Errorf("Expected %s, but got %s", expectedDuration, result)
	}

	// Test invalid duration string
	durationStr = "invalid"
	expectedDuration = 0
	result = ParseDuration(durationStr, expectedDuration)
	if result != expectedDuration {
		t.Errorf("Expected %s, but got %s", expectedDuration, result)
	}

	// Test empty duration string
	durationStr = ""
	expectedDuration = 0
	result = ParseDuration(durationStr, expectedDuration)
	if result != expectedDuration {
		t.Errorf("Expected %s, but got %s", expectedDuration, result)
	}
}

func TestSplitDateTime(t *testing.T) {
	tests := []struct {
		name         string
		dateTimeStr  string
		expectedDate string
		expectedTime string
	}{
		{
			name:         "valid_date_and_time",
			dateTimeStr:  "2024-11-11T12:10:00",
			expectedDate: "2024-11-11",
			expectedTime: "12:10",
		},
		{
			name:         "invalid_format_missing_T",
			dateTimeStr:  "2024-11-11 12:10:00",
			expectedDate: "",
			expectedTime: "",
		},
		{
			name:         "invalid_format_missing_time",
			dateTimeStr:  "2024-11-11T",
			expectedDate: "",
			expectedTime: "",
		},
		{
			name:         "invalid_format_missing_date",
			dateTimeStr:  "T12:10:00",
			expectedDate: "",
			expectedTime: "",
		},
		{
			name:         "invalid_format_missing_seconds_should_success",
			dateTimeStr:  "2024-11-11T12:10",
			expectedDate: "2024-11-11",
			expectedTime: "12:10",
		},
		{
			name:         "valid_date_and_time_with_different_time",
			dateTimeStr:  "2024-11-11T23:59:59",
			expectedDate: "2024-11-11",
			expectedTime: "23:59",
		},
		{
			name:         "valid_date_and_time_with_different_time_with_zone",
			dateTimeStr:  "2024-11-11T23:59:59Z",
			expectedDate: "2024-11-11",
			expectedTime: "23:59",
		},
		{
			name:         "valid_date_and_time_with_different_time_with_milis_and_zone",
			dateTimeStr:  "2025-04-23T09:24:07.303Z",
			expectedDate: "2025-04-23",
			expectedTime: "09:24",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date, time := SplitDateTime(tt.dateTimeStr)

			require.Equal(t, tt.expectedDate, date)
			require.Equal(t, tt.expectedTime, time)
		})
	}
}
