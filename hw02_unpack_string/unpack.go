package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/example/stringutil"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var isNumber bool
	var repeatCount int
	var output strings.Builder

	// unpack string from its end to beginning
	for i, char := range stringutil.Reverse(input) {
		digit, err := strconv.Atoi(string(char))
		if i == 0 && err != nil {
			output.WriteString(string(char))
		}

		if err == nil {
			if i == len(input)-1 {
				return "", ErrInvalidString
			}

			if isNumber {
				return "", ErrInvalidString
			}

			repeatCount = digit
			isNumber = true
			continue
		}

		output.WriteString(strings.Repeat(string(char), repeatCount))
		repeatCount = 1
		isNumber = false
	}

	return stringutil.Reverse(output.String()), nil
}
