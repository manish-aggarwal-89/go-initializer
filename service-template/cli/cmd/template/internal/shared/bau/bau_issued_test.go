package bau

import "testing"

func TestIssuanceHealthStatus_MarshalJSONx(t *testing.T) {
	// Test case 1: SUCCESS
	status := IssuanceHealthStatusSuccess
	expected := `"SUCCESS"`
	actual, err := status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 2: ISSUED_FAILED
	status = IssuanceHealthStatusIssuanceFailed
	expected = `"ISSUED_FAILED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 5: UNMAPPED
	status = IssuanceHealthStatusUnmapped
	expected = `"UNMAPPED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}
}

func TestIssuanceHealthStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    IssuanceHealthStatus
		want string
	}{
		{"SUCCESS", IssuanceHealthStatusSuccess, "SUCCESS"},
		{"ISSUED_FAILED", IssuanceHealthStatusIssuanceFailed, "ISSUED_FAILED"},
		{"UNMAPPED", 99, "UNMAPPED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("IssuedHealthStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssuanceLogHealthStatus_MarshalJSONx(t *testing.T) {
	// Test case 1: SUCCESS
	status := IssuanceLogHealthStatusSuccess
	expected := `"SUCCESS"`
	actual, err := status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 2: ISSUED_FAILED
	status = IssuanceLogHealthStatusIssuanceFailed
	expected = `"ISSUED_FAILED"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}

	// Test case 5: UNMAPPED
	status = IssuanceLogHealthStatusUnknownError
	expected = `"UNKNOWN_ERROR"`
	actual, err = status.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(actual) != expected {
		t.Errorf("expected %q, but got %q", expected, string(actual))
	}
}

func TestIssuanceLogHealthStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    IssuanceLogHealthStatus
		want string
	}{
		{"SUCCESS", IssuanceLogHealthStatusSuccess, "SUCCESS"},
		{"ISSUED_FAILED", IssuanceLogHealthStatusIssuanceFailed, "ISSUED_FAILED"},
		{"UNKNOWN_ERROR", 99, "UNKNOWN_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("IssuedLogHealthStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
