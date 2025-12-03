package bau

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tiket/TIX-FLIGHT-COMMON-MODEL-GO/model/bau"
)

func TestConstructAirlineProcessTimeDetail(t *testing.T) {
	// Test case 1: No analytic segments
	analyticSegments := []*bau.AnalyticSegment{}
	expected := map[string]int64{}
	actual := ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)

	// Test case 2: Single analytic segment
	analyticSegments = []*bau.AnalyticSegment{
		{
			Name:        "Segment 1",
			StartMillis: 0,
			EndMillis:   1000,
			IsParallel:  false,
		},
	}
	expected = map[string]int64{
		"Segment 1": 1000,
	}
	actual = ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)

	// Test case 3: Multiple analytic segments
	analyticSegments = []*bau.AnalyticSegment{
		{
			Name:        "Segment 1",
			StartMillis: 0,
			EndMillis:   1000,
			IsParallel:  false,
		},
		{
			Name:        "Segment 2",
			StartMillis: 1000,
			EndMillis:   2000,
			IsParallel:  false,
		},
		{
			Name:        "Segment 3",
			StartMillis: 2000,
			EndMillis:   3000,
			IsParallel:  false,
		},
	}
	expected = map[string]int64{
		"Segment 1": 1000,
		"Segment 2": 1000,
		"Segment 3": 1000,
	}
	actual = ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)

	// Test case 4: Multiple parallel analytic segments
	analyticSegments = []*bau.AnalyticSegment{
		{
			Name:        "Segment 1",
			StartMillis: 0,
			EndMillis:   1000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 2",
			StartMillis: 1000,
			EndMillis:   2000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 3",
			StartMillis: 2000,
			EndMillis:   3000,
			IsParallel:  true,
		},
	}
	expected = map[string]int64{
		"Segment 1": 1000,
		"Segment 2": 1000,
		"Segment 3": 1000,
	}
	actual = ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)

	// Test case 5: Multiple parallel analytic segments with different durations
	analyticSegments = []*bau.AnalyticSegment{
		{
			Name:        "Segment 1",
			StartMillis: 0,
			EndMillis:   1000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 2",
			StartMillis: 1000,
			EndMillis:   2000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 3",
			StartMillis: 2000,
			EndMillis:   3000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 4",
			StartMillis: 3000,
			EndMillis:   4000,
			IsParallel:  true,
		},
	}
	expected = map[string]int64{
		"Segment 1": 1000,
		"Segment 2": 1000,
		"Segment 3": 1000,
		"Segment 4": 1000,
	}
	actual = ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)

	// Test case 6: Multiple parallel analytic segments with the same name
	analyticSegments = []*bau.AnalyticSegment{
		{
			Name:        "Segment 1",
			StartMillis: 0,
			EndMillis:   1000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 1",
			StartMillis: 1000,
			EndMillis:   2000,
			IsParallel:  true,
		},
		{
			Name:        "Segment 1",
			StartMillis: 2000,
			EndMillis:   3000,
			IsParallel:  true,
		},
	}
	expected = map[string]int64{
		"Segment 1": 1000,
	}
	actual = ConstructAirlineProcessTimeDetail(analyticSegments)
	assert.Equal(t, expected, actual)
}

func TestSumMap(t *testing.T) {
	// Test case 1: empty map
	mapOfProcessingTime := make(map[string]int64)
	expectedResult := int64(0)
	result := SumMap(mapOfProcessingTime)
	assert.Equal(t, expectedResult, result)

	// Test case 2: non-empty map
	mapOfProcessingTime = map[string]int64{
		"A": 10,
		"B": 20,
		"C": 30,
	}
	expectedResult = 60
	result = SumMap(mapOfProcessingTime)
	assert.Equal(t, expectedResult, result)

	// Test case 3: map with negative values
	mapOfProcessingTime = map[string]int64{
		"A": -10,
		"B": -20,
		"C": -30,
	}
	expectedResult = -60
	result = SumMap(mapOfProcessingTime)
	assert.Equal(t, expectedResult, result)
}

func TestRemoveNumberSuffix(t *testing.T) {
	inps := []string{"", "abcd", "abc_09", "abc_3_32"}
	expected := []string{"", "abcd", "abc", "abc"}
	for i, inp := range inps {
		assert.Equal(t, expected[i], RemoveNumberSuffix(inp))
	}
}
