package shared

import (
	"fmt"
	"strings"
)

// ================================================================================================

// FormatString takes the targetFormat and args that are in key value pairs to format a string by
// replacing the keys in the format "{key}" with the value.
//
// Example:
//
//	// Output: "hello, world!"
//	shared.FormatString("hello, {name}!", "name", "world")
//	// Output: "the maximum number of cats is 10"
//	shared.FormatString("the maximum number of cats is {maxCats}", "maxCats", 10)
//
// Note: If the args are missing, odd, etc., it does not panic. It simply tries to format the
// string as best as it can.
func FormatString(targetFormat string, args ...any) string {
	// if there aren't any pairs, there's nothing to format
	numPairs := len(args) / 2
	if numPairs < 1 {
		return targetFormat
	}

	// format keys to have "{" and "}" surrounding the key name and to make sure the args are strings
	evenLength := numPairs * 2
	evenArgs := args[:evenLength]
	formattedArgs := make([]string, evenLength)
	for idx := 0; idx < len(evenArgs); idx = idx + 2 {
		key := StringOrEmpty(evenArgs[idx])
		value := StringOrEmpty(evenArgs[idx+1])
		formattedArgs[idx] = fmt.Sprintf("{%s}", key)
		formattedArgs[idx+1] = value
	}

	// use a replacer to replace our key value pairs
	replacer := strings.NewReplacer(formattedArgs...)
	return replacer.Replace(targetFormat)
}

// ------------------------------------------------------------------------------------------------

// StringOrEmpty returns an emptry string if target is nil, otherwise it uses "%v" to format
// target.
func StringOrEmpty(target any) string {
	if target == nil {
		return ""
	}

	return fmt.Sprintf("%v", target)
}
