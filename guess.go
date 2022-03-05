package gocuke

import (
	"regexp"
	"strings"
	"unicode"
)

func guessMethodSig(text string) methodSig {
	parts := strings.Split(strings.TrimSpace(text), " ")
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
				regexParts = append(regexParts, "("+decRegex.String()+")")
				continue
			}

			if intRegex.MatchString(part) {
				paramTypes = append(paramTypes, "int64")
				regexParts = append(regexParts, "("+intRegex.String()+")")
				continue
			}
		}

		nameParts = append(nameParts, part)
		regexParts = append(regexParts, part)
	}

	regex := regexp.MustCompile(strings.Join(regexParts, ` `))

	if len(nameParts) == 0 {
		return methodSig{name: "unknown", paramTypes: paramTypes, regex: regex}
	}

	var name string
	for i := 0; i < len(nameParts); i++ {
		n := toFirstUpperIdentifier(nameParts[i])
		if n == "" {
			continue
		}
		name = name + n
	}

	return methodSig{name: name, paramTypes: paramTypes, regex: regex}
}

func firstChar(x string) byte {
	return x[0]
}

func lastChar(x string) byte {
	return x[len(x)-1]
}

func toFirstUpperIdentifier(str string) string {
	runes := []rune(str)
	var res []rune
	isFirst := true
	for _, r := range runes {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			continue
		}

		if isFirst {
			res = append(res, unicode.ToUpper(r))
			isFirst = false
		} else {
			res = append(res, unicode.ToLower(r))
		}
	}
	return string(res)
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
