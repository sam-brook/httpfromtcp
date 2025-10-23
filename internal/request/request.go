package request

import (
	"bytes"
	"fmt"
	"io"
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

func getHttpVersion(b []byte) []byte {
	version_parts := bytes.Split(b, []byte("/"))
	if !bytes.Equal(version_parts[0], []byte("HTTP")) || len(version_parts) != 2 {
		return []byte{}
	}
	return version_parts[1]
}

func isUpperAlphabetic(b []byte) bool {
	for _, c := range b {
		if c < 65 || c > 90 {
			return false
		}
	}
	return true
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case initialised:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.State = done
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.State = done
		case done:
			break outer
		}
	}
	return read, nil
}

func newRequest() *Request {
	return &Request{
		State: initialised,
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buf := make([]byte, 1024)
	bufIdx := 0

	for request.State != done {
		n, err := reader.Read(buf[bufIdx:])
		if err != nil {
			return nil, err
		}

		bufIdx += n
		readN, err := request.parse(buf[:bufIdx+n])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[readN:bufIdx])
		bufIdx -= readN
	}
	return request, nil
}

var ERROR_MALFORMED_REQ_LINE = fmt.Errorf("Malformed request line")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	index := bytes.Index(b, SEPARATOR)
	if index == -1 {
		return nil, 0, nil
	}
	requestLine := b[:index]
	read := index + len(SEPARATOR)
	parts := bytes.Split(requestLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQ_LINE
	}

	method := parts[0]
	if !isUpperAlphabetic(method) {
		return nil, 0, fmt.Errorf("Non capital letter found in:\nMethod Name: %s\nRequest Line: %s", method, requestLine)
	}

	requestTarget := parts[1]

	version := getHttpVersion(parts[2])
	if !bytes.Equal(version, []byte("1.1")) {
		fmt.Printf("%s", version)
		return nil, 0, fmt.Errorf("We do not currently support HTTP versions other than 1.1: %s", version)
	}

	result := &RequestLine{
		Method:        string(method),
		RequestTarget: string(requestTarget),
		HttpVersion:   string(version),
	}
	return result, read, nil
}
