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
}

func getHttpVersion(s string) string {
	splitVersionStr := strings.Split(s, "/")
	if len(splitVersionStr) != 2 {
		return ""
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

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, fmt.Errorf("Error reading data: %v\n\t%s", err, b)
	}

	s := string(b)
	rqLine, err := parseRequestLine(s)
	if err != nil {
		return &Request{}, fmt.Errorf("Error parsing request line: %v\n\t%s", err, b)
	}

	result := &Request{
		RequestLine: *rqLine,
		Headers:     make(map[string]string),
		Body:        []byte(""),
	}
	return result, nil
}

func parseRequestLine(s string) (*RequestLine, error) {
	lines := strings.Split(s, "\r\n")
	requestLine := lines[0]
	splitRequestLine := strings.Split(requestLine, " ")

	method := splitRequestLine[0]
	fmt.Printf("\n")
	if !isUpperAlphabetic(method) {
		return &RequestLine{}, fmt.Errorf("Non capital letter found in:\nMethod Name: %s\nRequest Line: %s", method, requestLine)
	}

	requestTarget := splitRequestLine[1]

	version := getHttpVersion(splitRequestLine[2])
	if version != "1.1" {
		return &RequestLine{}, fmt.Errorf("We do not currently support HTTP versions other than 1.1 %s", version)
	}

	result := &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   version,
	}
	return result, nil
}
