package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func getLinesFromReader(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		var current_line strings.Builder
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				current_line.Write(data[:i])
				out <- current_line.String()
				current_line.Reset()
				data = data[i+1:]
			}
			current_line.Write(data)
		}
		defer close(out)
		defer f.Close()
	}()
	return out
}

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		net, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Connection has been accepted!\n")
		lines := getLinesFromReader(net)
		for line := range lines {
			fmt.Printf("%s\n", line)
		}
		fmt.Println("Connection to channel closed")
	}
}
