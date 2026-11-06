package shared_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

// TestSliceMap is a simple table-driven test for verifying that shared.SliceMap works as expected.
func TestSliceMap(testCtx *testing.T) {
	// Simple func that converts strings to ints, setting invalid strings to math.MaxInt
	// Used for testing SliceMap
	strToIntFn := func(valStr string, _ int) int {
		valInt, err := strconv.Atoi(valStr)
		if err != nil {
			valInt = math.MaxInt
		}
		return valInt
	}

	testCases := []struct {
		Name         string
		Data         []string
		MappingFn    func(string, int) int
		ExpectedData []int
	}{
		{
			Name:         "array with valid data",
			MappingFn:    strToIntFn,
			Data:         []string{"1", "2", "3"},
			ExpectedData: []int{1, 2, 3},
		},
		{
			Name:         "array with invalid data",
			MappingFn:    strToIntFn,
			Data:         []string{"1", "2", "aa"},
			ExpectedData: []int{1, 2, math.MaxInt},
		},
		{
			Name:         "empty input array",
			MappingFn:    strToIntFn,
			Data:         []string{},
			ExpectedData: nil,
		},
		{
			Name:         "nil input array",
			MappingFn:    strToIntFn,
			Data:         nil,
			ExpectedData: nil,
		},
		{
			Name:         "nil mapping func",
			MappingFn:    nil,
			Data:         []string{"1", "2", "3"},
			ExpectedData: nil,
		},
	}

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			mappedData := shared.SliceMap(testCase.Data, testCase.MappingFn)
			assert.Equal(subTestCtx, testCase.ExpectedData, mappedData)
		})
	}
}
