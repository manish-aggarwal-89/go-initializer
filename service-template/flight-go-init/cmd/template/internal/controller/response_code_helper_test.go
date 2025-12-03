package controller

import (
	"testing"
	"{{MODULE_NAME}}/internal/shared"

	"github.com/stretchr/testify/require"
)

func TestGetBaggageResponseCodeMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "success",
			input:    nil,
			expected: shared.SUCCESS_RESPONSE_CODE,
		},
		{
			name:     "ancillary_is_empty",
			input:    shared.ErrBaggageAncillaryIsEmpty,
			expected: shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE,
		},
		{
			name:     "itinerary_length_not_match",
			input:    shared.ErrItineraryLengthNotMatch,
			expected: shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE,
		},
		{
			name:     "itinerary_not_found",
			input:    shared.ErrItineraryNotFound,
			expected: shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE,
		},
		{
			name:     "baggage_response_is_empty",
			input:    shared.ErrBaggageResponseIsEmpty,
			expected: shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE,
		},
		{
			name:     "default_case",
			input:    shared.ErrRequestAirlineFailed,
			expected: shared.BAGGAGE_NOT_AVAILABLE_RESPONSE_CODE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, GetBaggageResponseCodeMapper(tt.input))
		})
	}
}

func TestBookingResponseCodeMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "success",
			input:    nil,
			expected: shared.SUCCESS_RESPONSE_CODE,
		},
		{
			name:     "no_seat",
			input:    shared.ErrNoSeat,
			expected: shared.NO_SEAT_RESPONSE_CODE,
		},
		{
			name:     "journey_not_found",
			input:    shared.ErrJourneyNotFound,
			expected: shared.NO_SEAT_RESPONSE_CODE,
		},
		{
			name:     "default_case",
			input:    shared.ErrRequestAirlineFailed,
			expected: shared.BOOKING_FAILED_RESPONSE_CODE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, BookingResponseCodeMapper(tt.input))
		})
	}
}
