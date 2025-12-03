package util

import (
	"fmt"
	"testing"
	"{{MODULE_NAME}}/internal/shared"

	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/common"

	"github.com/stretchr/testify/assert"
)

func TestLogRest(t *testing.T) {
	mr := common.MandatoryRequest{
		RequestId: "1234567890",
	}

	ctx := "TestContext"
	rqrs := "TestRQRS"
	data := "TestData"

	expected := fmt.Sprintf("[%s][%s] [%s] - %s - %s", shared.REST_IMPL, ctx, rqrs, ObjToJson(mr), ObjToJson(data))
	result := LogRest(mr, ctx, rqrs, data)

	if result != expected {
		t.Errorf("LogRest() = %s; want %s", result, expected)
	}
}

func TestLogService(t *testing.T) {
	mr := common.MandatoryRequest{
		RequestId: "12345",
	}

	ctx := "TestContext"
	rqrs := "TestRQRS"
	data := map[string]string{"key": "value"}

	expected := "[SERVICE_IMPL][TestContext] [TestRQRS] - [12345] - {\"key\":\"value\"}"
	result := LogService(mr, ctx, rqrs, data)

	if result != expected {
		t.Errorf("LogService was incorrect, got: %s, want: %s.", result, expected)
	}
}

func TestLogOutbound(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		mr       common.MandatoryRequest
		ctx      string
		rqrs     string
		data     string
		expected string
	}{
		{
			name: "Test Case 1",
			mr: common.MandatoryRequest{
				RequestId: "123",
			},
			ctx:      "TestContext",
			rqrs:     "TestRQRS",
			data:     "TestData",
			expected: "[OUTBOUND_IMPL][TestContext] [TestRQRS] - [123] - TestData",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LogOutbound(tt.mr, tt.ctx, tt.rqrs, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogInbound(t *testing.T) {
	// Test case 1: Valid input
	mr := common.MandatoryRequest{
		RequestId: "123",
	}
	ctx := "test-ctx"
	rqrs := "test-rqrs"
	data := map[string]string{"key": "value"}
	expected := fmt.Sprintf("[%s][%s] [%s] - %s - %s", shared.INBOUND_IMPL, ctx, rqrs, ObjToJson(mr), ObjToJson(data))
	result := LogInbound(mr, ctx, rqrs, data)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case 2: Invalid input
	mr = common.MandatoryRequest{}
	ctx = ""
	rqrs = ""
	data = nil
	expected = fmt.Sprintf("[%s][%s] [%s] - %s - %s", shared.INBOUND_IMPL, ctx, rqrs, ObjToJson(mr), ObjToJson(data))
	result = LogInbound(mr, ctx, rqrs, data)
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestLogScheduler(t *testing.T) {
	// Test case 1: Valid input
	result := LogScheduler("test_project", "test_message")
	expected := "[scheduler][test_project] - test_message"
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case 2: Invalid input (empty project name)
	result = LogScheduler("", "test_message")
	expected = "[scheduler][] - test_message"
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	// Test case 3: Invalid input (empty message)
	result = LogScheduler("test_project", "")
	expected = "[scheduler][test_project] - "
	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
