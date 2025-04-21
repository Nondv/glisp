package reader

import (
	"regexp"
	"strconv"
	"unicode"
)

type Symbol string

const (
	SymbolReference = "sym"
	IntegerReference = "int"
)

type Value struct {
	valueType string
	value     any
}

func Read(txt string) Value {
	runes := []rune(txt)
	token, _, err := nextToken(runes)
	panicIfErr(err, "!!!")

	match, err := regexp.MatchString("^-?\\d+$", token)
	panicIfErr(err, "Something wrong with integer regexp")
	if match {
		value, err := strconv.Atoi(token)
		panicIfErr(err, "Can't parse int!")
		return Value{IntegerReference, value}
	}

	return Value{SymbolReference, &token}
}


func panicIfErr(err error, msg string) {
	if err != nil {
		panic(msg)
	}
}

// (token, index where it ends, error)
func nextToken(runes []rune) (string, int, error) {
	i := 0
	for unicode.IsSpace(runes[i]) {
		i += 1
	}

	beginIndex := i
	for i < len(runes) && !unicode.IsSpace(runes[i]) {
		i += 1
	}

	return string(runes[beginIndex:i]), i - 1, nil
}
