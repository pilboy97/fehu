package cli

import (
	"strings"
	"unicode"
)

func TokenizeCommand(cmd string) []string {
	if len(cmd) == 0 {
		return []string{}
	}

	str := []rune(cmd)
	//convert string to unicode string
	token := make([]string, 0, len(str))
	// tokens from command string

	var inQuote = false
	// is the ith character in quote or not
	var start = 0
	// start of the current token

	for i := 0; i < len(str); i++ {
		if unicode.IsSpace(str[i]) {
			if inQuote {
				continue
			}

			if start != -1 {
				token = append(token, string(str[start:i]))
			}

			start = -1
		} else {
			if str[i] == '"' {
				inQuote = !inQuote
			}
			if start == -1 {
				start = i
			}
		}
	}

	if start != -1 {
		token = append(token, string(str[start:]))

		if inQuote {
			return nil
		}
	}

	return token
}
func ParseFlag(str string) []string {
	token := strings.SplitN(str, "=", 2)

	if strings.HasPrefix(token[0], "--") {
		token[0] = strings.TrimPrefix(token[0], "--")
	} else if strings.HasPrefix(token[0], "-") {
		token[0] = strings.TrimPrefix(token[0], "-")
	} else {
		return nil
	}

	return token
}
