package bau

import (
	"reflect"
	"testing"
)

func TestBookingHealthStatus_MarshalJSON(t *testing.T) {
	// Test success case
	var status BookingHealthStatus = BookingHealthStatusSuccess
	expectedJSON := []byte(`"SUCCESS"`)
	actualJSON, err := status.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("Expected JSON %v, but got %v", expectedJSON, actualJSON)
	}

	// Test booking failed case
	status = BookingHealthStatusBookingFailed
	expectedJSON = []byte(`"BOOKING_FAILED"`)
	actualJSON, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("Expected JSON %v, but got %v", expectedJSON, actualJSON)
	}

	// Test no seat case
	status = BookingHealthStatusNoSeat
	expectedJSON = []byte(`"NO_SEAT"`)
	actualJSON, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("Expected JSON %v, but got %v", expectedJSON, actualJSON)
	}

	// Test unmapped case
	status = 999
	expectedJSON = []byte(`"UNMAPPED"`)
	actualJSON, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("Expected JSON %v, but got %v", expectedJSON, actualJSON)
	}
}

func TestBookingHealthStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		input    BookingHealthStatus
		expected string
	}{
		{
			name:     "Success",
			input:    BookingHealthStatusSuccess,
			expected: "SUCCESS",
		},
		{
			name:     "BookingFailed",
			input:    BookingHealthStatusBookingFailed,
			expected: "BOOKING_FAILED",
		},
		{
			name:     "NoSeat",
			input:    BookingHealthStatusNoSeat,
			expected: "NO_SEAT",
		},
		{
			name:     "Unmapped",
			input:    4,
			expected: "UNMAPPED",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.String()
			if result != test.expected {
				t.Errorf("Expected %s, but got %s", test.expected, result)
			}
		})
	}
}
