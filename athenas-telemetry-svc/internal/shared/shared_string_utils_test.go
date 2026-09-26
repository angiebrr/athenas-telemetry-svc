package shared_test

import (
	"testing"

	"github.com/angiebrr/athenas-telemetry-svc/internal/shared"
	"github.com/stretchr/testify/assert"
)

// ================================================================================================

func TestFormatString(testCtx *testing.T) {
	testCases := []struct {
		Name           string
		InputFormatStr string
		KeyValues      []any
		ExpectedStr    string
	}{
		{
			Name:           "format string with 1 variable used once",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"name", "world"},
			ExpectedStr:    "hello, world!",
		},
		{
			Name:           "format string with 1 variable used many times",
			InputFormatStr: "Sailor {sailorScoutName} is the strongest because {sailorScoutName} never gives up",
			KeyValues:      []any{"sailorScoutName", "Mars"},
			ExpectedStr:    "Sailor Mars is the strongest because Mars never gives up",
		},
		{
			Name:           "format string with multiple variables",
			InputFormatStr: "my name is {name} and I'm from {place}, and I have a dog named {dogName}",
			KeyValues:      []any{"name", "Tommy", "place", "England", "dogName", "Sammy"},
			ExpectedStr:    "my name is Tommy and I'm from England, and I have a dog named Sammy",
		},
		{
			Name:           "format string with extra variables",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"name", "world", "sailorName", "Mars"},
			ExpectedStr:    "hello, world!",
		},
		{
			Name:           "format string doesn't have variables",
			InputFormatStr: "who is your favorite Sailor Scout?",
			KeyValues:      nil,
			ExpectedStr:    "who is your favorite Sailor Scout?",
		},
		{
			Name:           "format string has variables but none match",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"sailorName", "Mars"},
			ExpectedStr:    "hello, {name}!",
		},
		{
			Name:           "format string has odd key value pairs (no matching variable)",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"name"},
			ExpectedStr:    "hello, {name}!",
		},
		{
			Name:           "format string has odd key value pairs (matching variable)",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"name", "Rei San", "favoriteColor"},
			ExpectedStr:    "hello, Rei San!",
		},
		{
			Name:           "format string with 1 variable with {} surrounding value",
			InputFormatStr: "hello, {name}!",
			KeyValues:      []any{"name", "{world}"},
			ExpectedStr:    "hello, {world}!",
		},
		{
			Name:           "format string is empty",
			InputFormatStr: "",
			ExpectedStr:    "",
		},
		{
			Name:           "value of a key is the same name as another key",
			InputFormatStr: "{name}",
			KeyValues:      []any{"name", "{wow}", "wow", "???"},
			ExpectedStr:    "{wow}",
		},
		{
			Name:           "multiple of the same keys are defined",
			InputFormatStr: "hi, {name}!",
			KeyValues:      []any{"name", "Verity", "name", "Imposter"},
			ExpectedStr:    "hi, Verity!",
		},
		{
			Name:           "nil keys and values",
			InputFormatStr: "hi, {name}! are you {<nil>}?",
			KeyValues:      []any{"name", nil, nil, "nothingness"},
			ExpectedStr:    "hi, ! are you {<nil>}?",
		},
		{
			Name:           "values that are numbers",
			InputFormatStr: "I have {computerCount} computers",
			KeyValues:      []any{"computerCount", 10},
			ExpectedStr:    "I have 10 computers",
		},
	}

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			formattedStr := shared.FormatString(testCase.InputFormatStr, testCase.KeyValues...)
			got := formattedStr
			want := testCase.ExpectedStr
			assert.Equal(subTestCtx, want, got)
		})
	}
}

// ------------------------------------------------------------------------------------------------

func TestStringOrEmpty(testCtx *testing.T) {
	testCases := []struct {
		Name           string
		Input          any
		ExpectedOutput string
	}{
		{Name: "input string", Input: "hello", ExpectedOutput: "hello"},
		{Name: "nil input should be empty string", Input: nil, ExpectedOutput: ""},
		{Name: "input int", Input: 0, ExpectedOutput: "0"},
	}

	for _, testCase := range testCases {
		testCtx.Run(testCase.Name, func(subTestCtx *testing.T) {
			outputStr := shared.StringOrEmpty(testCase.Input)
			assert.Equal(subTestCtx, testCase.ExpectedOutput, outputStr)
		})
	}
}
