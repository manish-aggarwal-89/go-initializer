package bau

import (
	"testing"
)

func TestSearchHealthStatus_MarshalJSON(t *testing.T) {
	// Test case 1: SUCCESS
	status := SearchingHealthStatusSuccess
	expected := `"SUCCESS"`
	actual, err := status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 2: SEARCH_FAILED
	status = SearchingHealthStatusSearchFailed
	expected = `"SEARCH_FAILED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 3: FLIGHT_NOT_FOUND
	status = SearchingHealthStatusFlightNotFound
	expected = `"FLIGHT_NOT_FOUND"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 4: GET_FARE_FAILED
	status = SearchingHealthStatusGetFareFailed
	expected = `"GET_FARE_FAILED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 5: UNMAPPED
	status = SearchingHealthStatusUnmapped
	expected = `"UNMAPPED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 6: FLIGHT_NOT_FOUND
	status = SearchingHealthStatusInvalidRoute
	expected = `"INVALID_ROUTE"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}
}

func TestSearchHealthStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    SearchHealthStatus
		want string
	}{
		{"SUCCESS", SearchingHealthStatusSuccess, "SUCCESS"},
		{"SEARCH_FAILED", SearchingHealthStatusSearchFailed, "SEARCH_FAILED"},
		{"FLIGHT_NOT_FOUND", SearchingHealthStatusFlightNotFound, "FLIGHT_NOT_FOUND"},
		{"GET_FARE_FAILED", SearchingHealthStatusGetFareFailed, "GET_FARE_FAILED"},
		{"INVALID_ROUTE", SearchingHealthStatusInvalidRoute, "INVALID_ROUTE"},
		{"UNMAPPED", 99, "UNMAPPED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("SearchHealthStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
