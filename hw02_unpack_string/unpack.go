package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputStr string) (string, error) {
	if inputStr == "" {
		return "", nil
	}

	rn, runeSz := utf8.DecodeRuneInString(inputStr)
	if unicode.IsDigit(rn) {
		return "", ErrInvalidString
	}

	var outputStr strings.Builder
	inputStr = inputStr[runeSz:]
	outputStr.WriteString(string(rn))

	prevRn := rn
	for _, curRn := range inputStr {
		if unicode.IsDigit(prevRn) && unicode.IsDigit(curRn) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(curRn) {
			cnt, _ := strconv.Atoi(string(curRn))
			if cnt > 0 {
				outputStr.WriteString(strings.Repeat(string(prevRn), cnt-1))
			} else {
				suf := string(prevRn)
				tmpStr := outputStr.String()
				tmpStr = strings.TrimSuffix(tmpStr, suf)
				outputStr.Reset()
				outputStr.WriteString(tmpStr)
			}
		} else {
			outputStr.WriteString(string(curRn))
		}
		prevRn = curRn
	}
	res := outputStr.String()
	return res, nil
}
