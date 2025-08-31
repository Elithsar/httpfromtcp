package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	crlfIndex := -1
	for i := 0; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			crlfIndex = i
			break
		}
	}

	if crlfIndex == -1 {
		// No CRLF found, need more data
		return 0, false, nil
	}

	if crlfIndex == 0 {
		// Found empty line, end of headers
		return 2, true, nil // Consume the CRLF
	}

	// Parse the header line
	line := string(data[:crlfIndex])
	colonIndex := -1
	for i := 0; i < len(line); i++ {
		if line[i] == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex == -1 || colonIndex == 0 || colonIndex == len(line)-1 {
		return 0, false, fmt.Errorf("invalid header format")
	}

	if colonIndex > 0 && line[colonIndex-1] == ' ' {
		return 0, false, fmt.Errorf("invalid header format: space before colon")
	}

	key := strings.TrimSpace(line[:colonIndex])

	if strings.Contains(key, " ") {
		return 0, false, fmt.Errorf("invalid header format: space in key")
	}
	key = strings.ToLower(key)
	for _, ch := range key {
		if !(ch >= 'a' && ch <= 'z') &&
			!(ch >= '0' && ch <= '9') &&
			!strings.ContainsRune("!#$%&'*+-.^_`|~", ch) {
			return 0, false, fmt.Errorf("invalid character in header key: %q", ch)
		}
	}

	value := strings.TrimSpace(line[colonIndex+1:])

	if existingValue, exists := h[key]; exists {
		value = existingValue + ", " + value
	}
	h[key] = value

	return crlfIndex + 2, false, nil // +2 to consume the CRLF
}
