package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var twoDigitsInARow bool
	var previousChar string
	var output strings.Builder

	for i, char := range input {
		currentChar := string(char)
		repeatCount, err := strconv.Atoi(currentChar)
		if err != nil {
			if !twoDigitsInARow {
				output.WriteString(previousChar)
			}
			if i == len(input)-1 {
				output.WriteString(currentChar)
			}

			previousChar = currentChar
			twoDigitsInARow = false
			continue
		}

		if i == 0 {
			return "", ErrInvalidString
		}
		if twoDigitsInARow {
			return "", ErrInvalidString
		}

		output.WriteString(strings.Repeat(previousChar, repeatCount))
		twoDigitsInARow = true
	}

	return output.String(), nil
}
