package reader

import (
	"container/list"
	"errors"
	"regexp"
	"strconv"
	"unicode"

	. "nondv.io/glisp/types"
)

// Reads only 1 sexp
func Read(txt string) (*Value, error) {
	runes := []rune(txt)

	value, _, err := parseNext(runes)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// func printTokens(t *testing.T, runes []rune) {
// 	token, i, err := nextToken(runes)
// 	for err == nil && i+1 < len(runes) {
// 		t.Logf("%s\n", token)
// 		var iOffset int
// 		token, iOffset, err = nextToken(runes[i+1:])
// 		i += 1 + iOffset
// 	}
// 	if err == nil {
// 		t.Logf("%s\n", token)
// 	}
// }

func parseNext(runes []rune) (*Value, int, error) {
	token, endIndex, err := nextToken(runes)
	if err != nil {
		return nil, endIndex, err
	}

	if token == "(" {
		return parseList(runes)
	}

	value, err := tokenToValue(token)
	return value, endIndex, err
}

func parseList(runes []rune) (*Value, int, error) {
	token, i, err := nextToken(runes)
	if err != nil || token != "(" {
		panic("parseList didn't find an opening paren")
	}

	values := list.New()
	token, maybeClosingParenOffset, _ := nextToken(runes[i+1:])
	for token != ")" {
		val, iOffset, err := parseNext(runes[i+1:])
		values.PushFront(val)
		i += 1 + iOffset
		if err != nil {
			return nil, i, err
		}

		token, maybeClosingParenOffset, _ = nextToken(runes[i+1:])
	}

	result := BuildEmptyList()
	for e := values.Front(); e != nil; e = e.Next() {
		result = BuildCons(e.Value.(*Value), result)
	}

	return result, i + 1 + maybeClosingParenOffset, nil
}

func tokenToValue(token string) (*Value, error) {
	if token == "(" || token == ")" {
		return nil, errors.New("Can't convert to value")
	}

	isInteger, err := regexp.MatchString("^-?\\d+$", token)
	panicIfErr(err)
	if isInteger {
		value, err := strconv.Atoi(token)
		panicIfErr(err)
		return BuildInteger(value), nil
	}

	return BuildSymbol(token), nil
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// (token, index where it ends, error)
func nextToken(runes []rune) (string, int, error) {
	i := 0
	for unicode.IsSpace(runes[i]) {
		i += 1
	}

	if isParen(runes[i]) {
		return string(runes[i]), i, nil
	}

	beginIndex := i
	for i < len(runes) && !unicode.IsSpace(runes[i]) && !isParen(runes[i]) {
		i += 1
	}

	return string(runes[beginIndex:i]), i - 1, nil
}

func isParen(r rune) bool {
	return r == '(' || r == ')'
}
