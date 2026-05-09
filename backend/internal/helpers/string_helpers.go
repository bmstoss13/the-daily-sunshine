package helpers

import (
	"fmt"
	"strconv"
)

func StringToBase10Int32(str string) (int32, error) {
	//bitSize set to 32 to ensure fits into int32
	i64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("[string_helpers.go] StringToBase10Int32: An error occurred while converting string %v to int32: %w", str, err)
	}

	// cast resulting int64 to int32
	i32 := int32(i64)
	return i32, err
}
