package gocuke

import (
	"regexp"
	"strings"
	"unicode"
)

func guessMethodSig(text string) methodSig {
	parts := strings.Split(text, " ")
	var (
		nameParts     []string
		paramTypes    []string
		regexParts    []string
		inSingleQuote bool
		inDoubleQuote bool
	)
	for i := 0; i < len(parts); i++ {
		part := parts[i]

		if inSingleQuote {
			if lastChar(part) == '\'' {
				inSingleQuote = false
			}
			continue
		}

		if inDoubleQuote {
			if lastChar(part) == '"' {
				inDoubleQuote = false
			}
			continue
		}

		c := firstChar(part)
		switch c {
		case '\'':
			if lastChar(part) != '\'' {
				inSingleQuote = true
			}
			paramTypes = append(paramTypes, "string")
			regexParts = append(regexParts, `'([^']*)'`)
			continue
		case '"':
			if lastChar(part) != '"' {
				inDoubleQuote = true
			}
			paramTypes = append(paramTypes, "string")
			regexParts = append(regexParts, `"([^"]*)"`)
			continue
		default:
			if decRegex.MatchString(part) {
				paramTypes = append(paramTypes, "float64")
				regexParts = append(regexParts, decRegex.String())
				continue
			}

			if intRegex.MatchString(part) {
				paramTypes = append(paramTypes, "int64")
				regexParts = append(regexParts, intRegex.String())
				continue
			}
		}

		nameParts = append(nameParts, part)
	}

	if len(nameParts) == 0 {
		return methodSig{name: "unknown", paramTypes: paramTypes}
	}

	name := strings.ToLower(nameParts[0])
	for i := 1; i < len(nameParts); i++ {
		name = name + toFirstUpper(nameParts[i])
	}

	return methodSig{name: name, paramTypes: paramTypes}
}

func firstChar(x string) byte {
	return x[0]
}

func lastChar(x string) byte {
	return x[len(x)-1]
}

func toFirstUpper(str string) string {
	runes := []rune(str)
	isFirst := true
	for i, r := range runes {
		if isFirst {
			runes[i] = unicode.ToUpper(r)
			isFirst = false
		} else {
			runes[i] = unicode.ToLower(r)
		}
	}
	return string(runes)
}

type methodSig struct {
	name       string
	paramTypes []string
	regex      *regexp.Regexp
}

var decRegex = regexp.MustCompile(`\d+(\.\d+)`)
var intRegex = regexp.MustCompile(`\d+`)

//func formatMethodStub(step *messages.PickleStep)  {
//	sig := guessMethodSig(step.Text)
//	TODO
//}
