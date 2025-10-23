package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	State       State
}

type State int

const (
	initialised State = iota
	done
)

func getHttpVersion(s string) string {
	splitVersionStr := strings.Split(s, "/")
	if len(splitVersionStr) != 2 || splitVersionStr[0] != "HTTP" {
		return s
	}
	return splitVersionStr[1]
}

func isUpperAlphabetic(s string) bool {
	for _, c := range s {
		if c < 65 || c > 90 {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == initialised {

	} else if r.State == done {

	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, fmt.Errorf("Error reading data: %v\n\t%s", err, b)
	}

	s := string(b)
	rqLine, state, err := parseRequestLine(s)
	if err != nil {
		if state == 
		return &Request{}, fmt.Errorf("Error parsing request line: %v\n\t%s", err, b)
	}

	result := &Request{
		RequestLine: *rqLine,
		Headers:     make(map[string]string),
		Body:        []byte(""),
	}
	return result, nil
}

func parseRequestLine(s string) (*RequestLine, int, error) {
	lines := strings.Split(s, "\r\n")
	if len(lines) == 0 {
		return &RequestLine{}, int(initialised), nil
	}
	requestLine := lines[0]
	splitRequestLine := strings.Split(requestLine, " ")

	method := splitRequestLine[0]
	fmt.Printf("\n")
	if !isUpperAlphabetic(method) {
		return &RequestLine{}, int(done), fmt.Errorf("Non capital letter found in:\nMethod Name: %s\nRequest Line: %s", method, requestLine)
	}

	requestTarget := splitRequestLine[1]

	version := getHttpVersion(splitRequestLine[2])
	if version != "1.1" {
		return &RequestLine{}, int(done), fmt.Errorf("We do not currently support HTTP versions other than 1.1 %s", version)
	}

	result := &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}
	return result,int(done), nil
}
