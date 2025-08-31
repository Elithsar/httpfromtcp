package request

import (
	"errors"
	"httpfromtcp/internal/headers"
	"io"
	"strings"
)

const (
	StateInitialized = iota
	StateRequestStateParsingHeaders
	StateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       int // 0 = initialized, 1 = done
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const (
	ErrInvalidRequestLine = "invalid request line"
)
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	r := &Request{
		state: StateInitialized,
	}

	for r.state != StateDone {
		// Grow buffer if full
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err == io.EOF {
			// Si on atteint EOF mais que l'état n'est pas StateDone, c'est une erreur
			if r.state != StateDone {
				return nil, errors.New("unexpected EOF: headers not terminated")
			}
			break
		}
		if err != nil {
			return nil, err
		}
		readToIndex += n

		// Parse the buffer up to readToIndex
		parsedBytes, parseErr := r.parse(buf[:readToIndex])
		if parseErr != nil {
			return nil, parseErr
		}

		// Remove parsed bytes from buffer
		if parsedBytes > 0 {
			copy(buf, buf[parsedBytes:readToIndex])
			readToIndex -= parsedBytes
		}
	}

	return r, nil
}

func parseRequestLine(data string) (RequestLine, int, error) {
	idx := strings.Index(data, "\r\n")
	if idx == -1 {
		return RequestLine{}, 0, nil
	}

	line := data[:idx]
	parts := strings.Split(line, " ")

	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New(ErrInvalidRequestLine)
	}

	method, requestTarget, httpVersion := parts[0], parts[1], parts[2]

	if strings.ToUpper(method) != method {
		return RequestLine{}, 0, errors.New(ErrInvalidRequestLine)
	}

	if httpVersion != "HTTP/1.1" {
		return RequestLine{}, 0, errors.New(ErrInvalidRequestLine)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   strings.TrimPrefix(httpVersion, "HTTP/"),
	}, idx + 2, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalParsed := 0
	for {
		switch r.state {
		case StateInitialized:
			// Parse request line
			strData := string(data[totalParsed:])
			requestLine, n, err := parseRequestLine(strData)
			if err != nil {
				return totalParsed + n, err
			}
			if n == 0 {
				// Not enough data for request line
				return totalParsed, nil
			}
			r.RequestLine = requestLine
			r.state = StateRequestStateParsingHeaders
			r.Headers = headers.NewHeaders()
			totalParsed += n
			// Continue to parse headers in next state
		case StateRequestStateParsingHeaders:
			n, done, err := r.Headers.Parse(data[totalParsed:])
			if err != nil {
				return totalParsed + n, err
			}
			if n == 0 && !done {
				// Need more data for headers
				return totalParsed, nil
			}
			totalParsed += n
			if done {
				r.state = StateDone
				return totalParsed, nil
			}
			// Continue parsing headers if not done
		case StateDone:
			return totalParsed, nil
		}
	}
}
