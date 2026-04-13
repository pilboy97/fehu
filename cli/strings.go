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

	var inQuote = false // True if currently inside a double quote
	var parenDepth = 0  // Tracks the nesting level of parentheses
	var start = 0       // Start index of the current token

	for i := 0; i < len(str); i++ {
		char := str[i]

		// Update quote and parenthesis states
		if char == '"' {
			inQuote = !inQuote
		} else if char == '(' {
			parenDepth++
		} else if char == ')' {
			parenDepth--
		}

		if unicode.IsSpace(char) {
			// If we are inside quotes or parentheses, treat space as part of the token
			if inQuote || parenDepth > 0 {
				continue
			}

			// If not inside quotes/parentheses, and we have a token started, append it
			if start != -1 {
				token = append(token, string(str[start:i]))
			}
			start = -1 // Reset start, indicating we are between tokens
		} else {
			// If not a space, and no token has started yet, mark current position as start
			if start == -1 {
				start = i
			}
		}
	}

	// After the loop, append the last token if any
	if start != -1 {
		token = append(token, string(str[start:]))
	}

	// Check for unclosed quotes or unbalanced parentheses at the end of the command
	if inQuote || parenDepth != 0 {
		// If there's an unclosed quote or unbalanced parentheses, the command is malformed.
		// Returning nil is consistent with the original behavior for unclosed quotes.
		return nil
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
